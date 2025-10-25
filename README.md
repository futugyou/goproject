# Golang Service

Projects in this monorepo are either `frequently used` or `frequently updated`,
some of them are deployed in **Vercel**.

The Older projects have been moved to [goproject-archived](https://github.com/futugyou/goproject-archived).

## Project Description

- **algorithm**/**design-pattern**: Some `algorithm` questions and `design pattern`
I did in the early years are basically not updated at present.

- **alphavantage-server**: Use [Alphavantage](www.alphavantage.co/premium/) to obtain data,
and provide API interface to [React](https://github.com/futugyou/react-project),
mainly for `React` to display reports.

- **aws-web**: Interact with AWS, the most important of which is `AppConfig`,
which provides data for the [React](https://github.com/futugyou/react-project) `network topology diagram`.

- **container**: Use golang to write `container` functions`, I don’t remember the original link.

- **identity-center**: Simple identity center for user login in projects like React.
Planned to be reimplemented in `identity-server`, currently on hold.

- **identity-mongo**: Using `MongoDB` to store `identity-center`.

- **infr-project**: The original `infr project`, which is a `DDD` project.
Currently undergoing reconstruction, it will only be used for **vercel** deployment in the future.

- **k8sbuilder**: `k8sbuilder` example, not updated yet.

- **openai-web**: Previously used to test the OpenAI SDK,
and has not been updated since switching to the `official SDK`.

- **tour**: For `go generate`, providing several simple functions.

- **other**: SDK or 'phonetic translation'.

## CircleCI/Sync Service/Dependabot Status

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/futugyou/goproject/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/futugyou/goproject/tree/master)
[![Alphavantage Data](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml)
[![SyncData](https://github.com/futugyou/goproject/actions/workflows/syncdata.yml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/syncdata.yml)
[![Dependabot](https://github.com/futugyou/goproject/actions/workflows/dependabot-auto.yml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/dependabot-auto.yml)

## CircleCI Insights

<details>

<summary style="font-size: 18px; font-weight: bold;">Infr Project Service</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/infr-project/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/infr-project/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>

<details>

<summary style="font-size: 18px; font-weight: bold;">AWS Service</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/aws-web/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/aws-web/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>

<details>

<summary style="font-size: 18px; font-weight: bold;">Identity Service</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/identity-center/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/identity-center/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>

<details>

<summary style="font-size: 18px; font-weight: bold;">OpenAI Service</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/openai-web/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/openai-web/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>

<details>

<summary style="font-size: 18px; font-weight: bold;">Alphavantage Service</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/alphavantage-server/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/alphavantage-server/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>

## GitHub Actions: Vercel Deployment Status

<details>

<summary>Vercel Deployment Status</summary>

> [!CAUTION]
> Due to a problem with `Vercel Cli`, the following action is `no longer available`.

[![AWS Preview](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-preview.yaml)
[![AWS Production](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/aws-vercel-production.yaml)

[![Identity Preview](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-preview.yaml)
[![Identity Production](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/identity-vercel-production.yaml)

[![OpenAI Preview](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-preview.yaml)
[![OpenAI Production](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/openAI-vercel-production.yaml)

[![Infr Project Preview](https://github.com/futugyou/goproject/actions/workflows/project-vercel-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/project-vercel-preview.yaml)
[![Infr Project Production](https://github.com/futugyou/goproject/actions/workflows/project-vercel-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/project-vercel-production.yaml)

[![Alphavantage Preview](https://github.com/futugyou/goproject/actions/workflows/alphavantage-preview.yaml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-preview.yaml)
[![Alphavantage Production](https://github.com/futugyou/goproject/actions/workflows/alphavantage-production.yaml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-production.yaml)

</details>
