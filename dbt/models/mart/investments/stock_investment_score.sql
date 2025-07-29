{{ config(materialized='table') }}

SELECT
  symbol,
  company_name,
  ROUND((roe_ttm + revenue_5y_growth + operating_margin) / 3, 2) AS investment_score,
  pe_ratio,
  dividend_yield
FROM {{ ref('stg_stock_scorecards') }}
WHERE roe_ttm IS NOT NULL AND revenue_5y_growth IS NOT NULL AND operating_margin IS NOT NULL
