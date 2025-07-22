from airflow import DAG
from airflow.providers.postgres.hooks.postgres import PostgresHook
from airflow.operators.python import PythonOperator
from airflow.exceptions import AirflowFailException
from datetime import datetime, timedelta, time as dtime
import pandas as pd
import talib
import requests
from datetime import timedelta
from pytz import timezone

default_args = {
    "start_date": datetime(2025, 7, 7),
    "email_on_failure": True,
    "email_on_retry": False,
    "retries": 1,
    "retry_delay": timedelta(minutes=5),
}

API_BASE_URL = "http://localhost:8081"
pg = PostgresHook(postgres_conn_id="postgres_local")
conn = pg.get_conn()
cursor = conn.cursor()


def get_last_processed_batch():
    cursor.execute(
        "SELECT last_run_batch_idx FROM airflow_progress_tracker WHERE airflow_dag = 'extract_intraday_data' ORDER BY created_at DESC LIMIT 1"
    )
    last_run_batch_idx = cursor.fetchone()
    last_run_batch_idx = last_run_batch_idx[0] if last_run_batch_idx else ""
    batch_id = 0

    if last_run_batch_idx != "":
        batch_id = last_run_batch_idx + 1

    def fetch_symbols(batch_id):
        cursor.execute(
            "SELECT symbol FROM stock_symbols WHERE batch_id = %s", (batch_id,)
        )
        return [r[0] for r in cursor.fetchall()]

    symbols = fetch_symbols(batch_id)

    if not symbols:
        batch_id = 0
        symbols = fetch_symbols(batch_id)

    return symbols, batch_id


def extract_stock_data(symbols, batch_id):
    """Extract stock data for a multiple symbols"""
    url = f"{API_BASE_URL}/api/extract/batch"
    eastern = timezone("US/Eastern")

    now = datetime.now(eastern)

    trading_start = eastern.localize(datetime.combine(now.date(), dtime(9, 30)))
    if now.time() > dtime(16, 0):
        now = eastern.localize(datetime.combine(now.date(), dtime(16, 0)))
    else:
        now = datetime.now(eastern)

    payload = {
        "symbols": symbols,
        "from": trading_start.isoformat(),
        "to": now.isoformat(),
    }

    try:
        response = requests.post(url, json=payload, timeout=300)
        response.raise_for_status()

        cursor.execute(
            "INSERT INTO airflow_progress_tracker (last_run_batch_idx, airflow_dag) VALUES (%s, %s)",
            (batch_id, "extract_intraday_data"),
        )
        conn.commit()

        print(f"Successfully extracted data for {symbols}")
        return True
    except Exception as e:
        print(f"Failed to extract data for {symbols}: {str(e)}")
        return False


def process_symbol_batch():
    symbols, batch_id = get_last_processed_batch()
    extract_intraday = extract_stock_data(symbols, batch_id)

    if not extract_intraday:
        print("Failed Extraction")
        return AirflowFailException(f"FAILED TO EXTRACT INTRADAY FOR {symbols}")


def get_last_processed_batch_metric_calc():
    cursor.execute(
        "SELECT last_run_batch_idx FROM airflow_progress_tracker WHERE airflow_dag = 'extract_intraday_data' ORDER BY created_at DESC LIMIT 1"
    )
    last_run_batch_idx = cursor.fetchone()
    last_processed_batch = last_run_batch_idx[0]

    cursor.execute(
        "SELECT * FROM stock_symbols WHERE batch_id = %s", (last_processed_batch,)
    )
    return [r[0] for r in cursor.fetchall()]


def calculate_indicators(symbol):
    pg = PostgresHook(postgres_conn_id="postgres_local")
    df = pg.get_pandas_df(
        f"""
        SELECT ip.*
        FROM stocks_intraday ip
        LEFT JOIN stocks_intraday_indicators ii
            ON ip.symbol = ii.symbol AND ip.date = ii.date
        WHERE ii.symbol IS NULL
        AND ip.date >= NOW() - interval '5 days'
        AND ip.symbol = '{symbol}'
        ORDER BY ip.date;
    """
    )

    results = []
    for symbol, group in df.groupby("symbol"):
        group = group.sort_values("date").reset_index(drop=True)

        open = group["open"].astype(float)
        high = group["high"].astype(float)
        low = group["low"].astype(float)
        close = group["close"].astype(float)
        volume = group["volume"].astype(float)

        # Indicators
        price_change_pct = (close - open) / open * 100
        avg_volume_50 = volume.rolling(window=50).mean()
        rvol = volume / avg_volume_50
        volume_spike = volume > 2 * avg_volume_50

        rsi = talib.RSI(close, timeperiod=78)
        atr = talib.ATR(high, low, close, timeperiod=78)
        obv = talib.OBV(close, volume)

        group["symbol"] = symbol
        group["price_change_pct"] = price_change_pct
        group["rvol"] = rvol
        group["volume_spike"] = volume_spike
        group["rsi"] = rsi
        group["atr"] = atr
        group["obv"] = obv

        valid_rows = group[
            [
                "symbol",
                "date",
                "rvol",
                "price_change_pct",
                "rsi",
                "volume_spike",
                "atr",
                "obv",
            ]
        ].dropna()
        valid_rows["volume_spike"] = valid_rows["volume_spike"].astype(bool)
        results += valid_rows.to_records(index=False).tolist()

    insert_query = """
        INSERT INTO intraday_indicators (
            symbol, date, rvol, price_change_pct, rsi, volume_spike, atr, obv
        ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
        ON CONFLICT (symbol, date) DO NOTHING
    """
    cursor.executemany(insert_query, results)


def calculate_metric_batch():
    symbols = get_last_processed_batch_metric_calc()
    for symbol in symbols:
        calculate_indicators(symbol)


with DAG(
    "extract_intraday_calculate_metrics",
    default_args=default_args,
    description="Extract intraday data and then calculate intraday metrics from intraday data that has been pulled",
    schedule=timedelta(minutes=5),
    catchup=False,
    tags=["intraday", "metric-calculation"],
) as dag:

    process_batch = PythonOperator(
        task_id="process_symbol_batch",
        python_callable=process_symbol_batch,
    )

    calculate_intraday_metrics = PythonOperator(
        task_id="calculate_intraday_metric_batch",
        python_callable=calculate_metric_batch,
    )

    process_batch >> calculate_intraday_metrics
