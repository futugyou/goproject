# Projectarium Service

## 1. Overview

This project is used to showcase the **`Projectarium`**, which contains 5 microservices and is deployed on the **Vercel** platform.

## 2. Description

- **`platformservice`**: This service is used to manage third-party platforms for **`Projectarium`** dependencies, such as GitHub, Vercel, and CircleCI.

- **`projectservice`**: This service is used to showcase **`Projectarium`**.

- **`resourcequeryservice`**: This project showcases the documentation resources for the **`Projectarium`**, specifically the `query` part of `CQRS`.

- **`resourceservice`**: This project is used to manipulate the documentation resources for the **`Projectarium`**, specifically the `command` part of `CQRS`.

- **`vaultservice`**: This service is used to handle the confidentiality required by **`Projectarium`**.

## 3. Deployment

- **Platform**: Vercel

- **CI/CD**: CircleCI for builds/tests

- **Monitoring & Status**: CircleCI Insights

## 4. Roadmap

- [ ] Use the `token exchange` function of the **`identity-center`**:.

- [ ] Processing `webhook` data from third-party platforms.

- [ ] **Monitoring**: Front-end projects are already connected to Grafana Cloud and Honeycomb.io; back-end integration pending.

## 5. Status Badges

<details open>

<summary style="font-size: 18px; font-weight: bold;">Projects Service</summary>

[![CircleCI](https://dl.circleci.com/insights-snapshot/gh/futugyou/goproject/master/infr-project/badge.svg?window=30d)](https://app.circleci.com/insights/github/futugyou/goproject/workflows/infr-project/overview?branch=master&reporting-window=last-30-days&insights-snapshot=true)

</details>
