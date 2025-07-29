{{ config(materialized='view') }}

SELECT
  symbol,
  fiscal_date,
  total_revenue,
  net_income
FROM {{ source('raw', 'income_statements') }}
