# FinStack – My Personal Finance + Stock Data Platform

This is a personal project I've been building to combine a bunch of financial tools and data into one place. It started as a way for me to sharpen my data engineering skills and try out some new tools and stacks — but it's also been super useful for managing my own investments and spending.

---

## 💡 What's this all about?

Basically, I wanted a one-stop app that could:
- Track my **bank transactions**
- Pull in **stock data** from various APIs
- Analyze both short-term signals and long-term fundamentals
- Let me mess around with real-world data engineering pipelines

I built a **Golang API** that pulls stock data from **Alpha Vantage**, **Finnhub**, and **Polygon**. Since most of these APIs have tight rate limits on the free tier, I had to design it to do **batch extraction** — so it can pull data over time without getting throttled.

On the personal finance side, I hooked into **Rows** (which uses Plaid in the background) to get access to my bank transaction history. I couldn’t connect directly to Plaid since that’s mostly for companies that have gone through audit processes, so this was a good workaround. Rows exports everything into Excel, and my backend reads that and loads it into my **Postgres RDS database**.

---

## 🔄 Automating it with Airflow

To automate all this, I used **Airflow**. It helps with:
- Regularly calling the stock APIs
- Tracking which data has already been pulled
- Calculating **metrics and indicators** (like RSI, volume spikes, long-term scorecards, etc.)

I'm hosting Airflow on **Astronomer**, but that's still a work in progress.

---

## 🛠 How it's deployed

All the deployments are done with **Terraform** — it gives me control over the different moving parts and makes it easy to spin things up or tear them down when needed.

Right now:
- The **Golang API** runs on an EC2 instance (public-facing)
- The **Postgres RDS** is private (only the EC2 can talk to it)
- Airbyte (also on EC2) syncs data from Postgres into **BigQuery**

---

## 📊 Analytics & DBT

I’m also playing around with **DBT** and **BigQuery** to do more in-depth analytics — like:
- What are my most frequent personal expenses?
- Which stocks are trending or undervalued based on custom scorecards?

Honestly, BigQuery is overkill for my scale, but I wanted to try it out in case I ever scaled this up

---

## 🖼 Frontend (WIP)

There's a **React frontend** in the works to show:
- Stock insights
- Portfolio breakdowns
- Personal finance dashboards

Still work in progress

---

## 💭 Why I built this

Mainly to learn:
- How to build reliable data pipelines
- How to connect different services together (cloud, open-source, SaaS)
- How to work with tools like **Airflow**, **Astronomer**, **Terraform**, **DBT**, **Airbyte**, etc.

It’s also just something I actually find useful in my daily life.

---

## 🧪 Stack Summary

| Part          | Tech |
|---------------|------|
| Backend       | Golang |
| Scheduler     | Airflow (Astronomer) |
| Data Movement | Airbyte |
| Analytics     | BigQuery + DBT |
| Infra         | Terraform (AWS + GCP) |
| DB            | PostgreSQL (RDS) |
| Frontend      | React (in progress) |

---

