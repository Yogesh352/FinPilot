CREATE TABLE IF NOT EXISTS airflow_progress_tracker (
    id SERIAL PRIMARY KEY,
    last_run_symbol VARCHAR(10) UNIQUE,
    airflow_dag VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS airflow_dag_idx ON airflow_progress_tracker(airflow_dag);

CREATE OR REPLACE FUNCTION update_airflow_progress_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_airflow_progress_updated_at 
    BEFORE UPDATE ON airflow_progress_tracker 
    FOR EACH ROW 
    EXECUTE FUNCTION update_airflow_progress_updated_at_column();