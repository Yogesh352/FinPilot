CREATE TABLE IF NOT EXISTS airflow_progress_tracker (
    id SERIAL PRIMARY KEY,
    last_run_batch_idx INTEGER,
    airflow_dag VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS airflow_dag_idx ON airflow_progress_tracker(airflow_dag);
CREATE INDEX IF NOT EXISTS airflow_batch_idx ON airflow_progress_tracker(last_run_batch_idx);