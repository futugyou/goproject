# alphavantage server

[![Alphavantage Data](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml)

## 1. Project Overview

Use the scheduler capability of GitHub Action to periodically obtain data from the AlphaVantage server for the [react project](https://github.com/futugyou/react-project) to display charts.

## 2. Scheduler Description

This service will be called by GitHub Actions at 23:30 UTC every day and retrieve data from [alphavantage](www.alphavantage.co).

## 3. Design

![1](./doc/images/arch.drawio.png)

<details open>

<summary style="font-size: 18px; font-weight: bold;">4. Alphavantage Service CircleCI Insights</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/alphavantage-server/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/alphavantage-server/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>

## 5. Miscellaneous

```shell
go mod edit -replace github.com/futugyou/alphavantage=../alphavantage

go mod edit -replace github.com/futugyou/alphavantage=github.com/futugyou/goproject/alphavantage@master
go mod tidy
```
