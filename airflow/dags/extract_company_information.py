from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow.providers.postgres.hooks.postgres import PostgresHook
import requests
from datetime import datetime, timedelta

default_args = {
    "owner": "airflow",
    "retries": 1,
    "retry_delay": timedelta(minutes=2),
}


def extract_next_batch():
    pg = PostgresHook(postgres_conn_id="postgres_local")
    conn = pg.get_conn()
    cursor = conn.cursor()

    cursor.execute(
        "SELECT last_run_batch_idx FROM airflow_progress_tracker WHERE airflow_dag = 'extract_company_profile_batch' ORDER BY created_at DESC LIMIT 1"
    )
    last_run_batch_idx = cursor.fetchone()
    last_run_batch_idx = last_run_batch_idx[0] if last_run_batch_idx else ""
    batch_id = 0

    if last_run_batch_idx != "":
        batch_id = last_run_batch_idx + 1

    # Step 2: Get the next 5 symbols
    cursor.execute(
        """
    SELECT symbol FROM stocks_metadata
    WHERE batch_id = %s
    """,
        (batch_id,),  # comma makes it a tuple
    )
    symbols = [row[0] for row in cursor.fetchall()]

    if not symbols:
        print("All symbols processed.")
        return

    # Call batch company extraction api
    print(f"Extracting for: {symbols}")
    response = requests.post(
        "http://localhost:8081/api/extract/companyprofile", json={"symbols": symbols}
    )
    if response.status_code == 200:
        print(f"Success for {symbols}")
        cursor.execute(
            "INSERT INTO airflow_progress_tracker (last_run_batch_idx, airflow_dag) VALUES (%s, %s)",
            (batch_id, "extract_company_profile_batch"),
        )
        conn.commit()
        print(f"Progress updated to {batch_id}")
    else:
        print(f"Failed for {symbols}: {response.status_code} {response.text}")


with DAG(
    dag_id="extract_company_profile_batch",
    default_args=default_args,
    schedule="*/2 * * * *",
    start_date=datetime(2023, 1, 1),
    catchup=False,
    max_active_runs=1,
    tags=["stocks", "api", "partitioned"],
) as dag:

    extract_batch = PythonOperator(
        task_id="extract_company_profile_batch", python_callable=extract_next_batch
    )
