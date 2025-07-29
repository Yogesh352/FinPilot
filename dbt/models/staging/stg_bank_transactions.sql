{{ config(materialized='view') }}

SELECT
  CAST(date AS DATE) AS transaction_date,
  amount,
  currency,
  description,
  category,
  bank,
  account
FROM {{ source('raw', 'bank_transactions') }}
