{{ config(materialized='table') }}

SELECT
  transaction_date,
  category,
  SUM(amount) AS total_spent
FROM {{ ref('stg_bank_transactions') }}
GROUP BY 1, 2
