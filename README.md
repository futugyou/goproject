# Golang Service

## 1. Project Overview

This monorepo contains projects that are either **frequently used** or **frequently updated**, some of which are deployed on **Vercel**.

Older projects have been moved to [goproject-archived](https://github.com/futugyou/goproject-archived).

## 2. Project Description

### 2.1 Main Projects

- **`alphavantage-server`**: Uses [Alphavantage](https://www.alphavantage.co/premium/) to obtain data and provides an API interface for [React](https://github.com/futugyou/react-project), mainly for displaying reports.

- **`aws-web`**: Interacts with AWS, primarily using `Aws Config` to provide data for the [React](https://github.com/futugyou/react-project) **network topology diagram**.

- **`identity-center`**: A simple identity center for user login in projects like React. Planned to be reimplemented as `identity-server` (currently on hold).

- **`infr-project`**: The original `infr-project` based on `DDD`. Currently being restructured; in the future, it will be used only for **Vercel** deployment.

- **`openai-web`**: Previously used to test the OpenAI SDK; no updates have been made since switching to the official SDK.

### 2.2 Other Projects

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

## 5. CircleCI & Github Actions Status

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/futugyou/goproject/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/futugyou/goproject/tree/master)
[![Alphavantage Data](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/alphavantage-data.yml)
[![SyncData](https://github.com/futugyou/goproject/actions/workflows/syncdata.yml/badge.svg?branch=master)](https://github.com/futugyou/goproject/actions/workflows/syncdata.yml)
[![Dependabot](https://github.com/futugyou/goproject/actions/workflows/dependabot-auto.yml/badge.svg)](https://github.com/futugyou/goproject/actions/workflows/dependabot-auto.yml)
