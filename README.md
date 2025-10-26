# Golang Service

## 1. Project Overview

This monorepo contains projects that are either **frequently used** or **frequently updated**, some of which are deployed on **Vercel**.

Older projects have been moved to [goproject-archived](https://github.com/futugyou/goproject-archived).

## 2. Project Description

### 2.1 Main Projects

- **`alphavantage-server`**: Uses [Alphavantage](https://www.alphavantage.co/premium/) to obtain data and provides an API interface for [React](https://github.com/futugyou/react-project), mainly for displaying reports.

- **`aws-web`**: Interacts with AWS, primarily using `AppConfig` to provide data for the [React](https://github.com/futugyou/react-project) **network topology diagram**.

- **`identity-center`**: A simple identity center for user login in projects like React. Planned to be reimplemented as `identity-server` (currently on hold).

- **`infr-project`**: The original `infr-project` based on `DDD`. Currently being restructured; in the future, it will be used only for **Vercel** deployment.

- **`openai-web`**: Previously used to test the OpenAI SDK; no updates have been made since switching to the official SDK.

### 2.2 Other Projects

- **`algorithm` / `design-pattern`**: Early works on `algorithm` problems and `design patterns`; mostly not updated currently.

- **`container`**: A `container` written in Go; the original link is not remembered.

- **`identity-mongo`**: Stores data for `identity-center` using **MongoDB**.

- **`k8sbuilder`**: Example project for `k8sbuilder`; not updated yet.

- **`tour`**: Provides several simple functions using `go generate`.

- **`other`**: Miscellaneous SDKs..

## 3. Deployment Architecture

- **Platform**: Vercel

- **CI/CD**: CircleCI for builds/tests, GitHub Actions for deployment (**currently non-functional**)

- **Monitoring & Status**: CircleCI Insights and GitHub Actions badges

## 4. Roadmap

- [ ] **`identity-center` → `identity-server`**: Reimplement to improve functionality.

- [ ] **`infr-project`**: Restructure for both unified and individual deployments.

- [ ] **Monitoring**: Front-end projects are already connected to Grafana Cloud and Honeycomb.io; back-end integration pending.

## 5. CircleCI/Sync Service/Dependabot Status

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/futugyou/goproject/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/futugyou/goproject/tree/master)
[![Alphavantage Data](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml)
[![SyncData](https://github.com/futugyou/goproject/actions/workflows/syncdata.yml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/syncdata.yml)
[![Dependabot](https://github.com/futugyou/goproject/actions/workflows/dependabot-auto.yml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/dependabot-auto.yml)

## 6. CircleCI Insights

<details open>

<summary style="font-size: 18px; font-weight: bold;">Infr Project Service</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/infr-project/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/infr-project/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>

<details open>

<summary style="font-size: 18px; font-weight: bold;">AWS Service</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/aws-web/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/aws-web/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>

<details open>

<summary style="font-size: 18px; font-weight: bold;">Identity Service</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/identity-center/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/identity-center/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>

<details open>

<summary style="font-size: 18px; font-weight: bold;">OpenAI Service</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/openai-web/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/openai-web/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>

<details open>

<summary style="font-size: 18px; font-weight: bold;">Alphavantage Service</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/alphavantage-server/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/alphavantage-server/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>

## 7. GitHub Actions: Vercel Deployment Status (`Archived`)

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
