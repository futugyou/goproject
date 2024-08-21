# Vercel Golang Service

## Dependabot

[![Dependabot](https://github.com/futugyou/goproject/actions/workflows/dependabot-auto.yml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/dependabot-auto.yml)

## Circle CI

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/futugyou/goproject/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/futugyou/goproject/tree/master)

## SyncData

[![SyncData](https://github.com/futugyou/goproject/actions/workflows/syncdata.yml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/syncdata.yml)

## AWS Service

[![AWS Preview](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-preview.yaml)
[![AWS Production](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-production.yaml)

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/aws-web/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/aws-web/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

## Identity Service

[![Identity Preview](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-preview.yaml)
[![Identity Production](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-production.yaml)

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/identity-center/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/identity-center/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

## OpenAI Service

[![OpenAI Preview](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-preview.yaml)
[![OpenAI Production](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-production.yaml)

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/openai-web/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/openai-web/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

## Infr Project Service

[![Infr Project Preview](https://github.com/futugyou/goproject/actions/workflows/project-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/project-vercel-preview.yaml)
[![Infr Project Production](https://github.com/futugyou/goproject/actions/workflows/project-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/project-vercel-production.yaml)

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/infr-project/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/infr-project/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

## Alphavantage Data

[![Alphavantage Data](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml)

## Alphavantage Service

[![Alphavantage Preview](https://github.com/futugyou/goproject/actions/workflows/alphavantage-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-preview.yaml)
[![Alphavantage Production](https://github.com/futugyou/goproject/actions/workflows/alphavantage-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-production.yaml)

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/alphavantage-server/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/alphavantage-server/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

## pin go version when go mod tidy

```golang
// set GOTOOLCHAIN except auto, eg. local path go1.20
go env -w  GOTOOLCHAIN=local
go mod tidy -go=1.20 -compat=1.20
```
