# Vercel Golang Service

## Dependabot

[![Dependabot](https://github.com/futugyou/goproject/actions/workflows/dependabot-auto.yml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/dependabot-auto.yml)

## SyncData

[![SyncData](https://github.com/futugyou/goproject/actions/workflows/syncdata.yml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/syncdata.yml)

## AWS Service

[![AWS Preview](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-preview.yaml)
[![AWS Production](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-production.yaml)

## Identity Service

[![Identity Preview](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-preview.yaml)
[![Identity Production](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-production.yaml)

## OpenAI Service

[![OpenAI Preview](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-preview.yaml)
[![OpenAI Production](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-production.yaml)

## Infr Project Service

[![Infr Project Preview](https://github.com/futugyou/goproject/actions/workflows/project-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/project-vercel-preview.yaml)
[![Infr Project Production](https://github.com/futugyou/goproject/actions/workflows/project-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/project-vercel-production.yaml)

## Alphavantage Data

[![Alphavantage Data](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml)

## Alphavantage Service

[![Alphavantage Preview](https://github.com/futugyou/goproject/actions/workflows/alphavantage-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-preview.yaml)
[![Alphavantage Production](https://github.com/futugyou/goproject/actions/workflows/alphavantage-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-production.yaml)

## pin go version when go mod tidy

```golang
// set GOTOOLCHAIN except auto, eg. local path go1.20
go env -w  GOTOOLCHAIN=local
go mod tidy -go=1.20 -compat=1.20
```
