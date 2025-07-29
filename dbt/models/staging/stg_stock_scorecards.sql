{{ config(materialized='view') }}

SELECT
  symbol,
  company_name,
  pe_ratio,
  peg_ratio,
  price_to_book,
  roe_ttm,
  revenue_5y_growth,
  operating_margin,
  profit_margin,
  dividend_yield,
  beta
FROM {{ source('raw', 'stock_scorecards') }}
