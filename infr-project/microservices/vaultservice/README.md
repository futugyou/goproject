# Vault Service

## 1. Overview

This project, part of **`Projectarium`**, is used to store `vaults`.

## 2. Description

Use `AES in CTR mode` to encrypt/decrypt the `vaults`.
In addition to local database storage, this service provides three additional storage methods:

1. **Azure**: use azure vault
2. **Aws**: use aws Simple Systems Manager (SSM)
3. **Hashicorp**: use hashicorp vault

## 3. Deployment Architecture

- **Vercel**: When deployed on the `Vercel`, this service and other `Projectarium` services will be deployed as a unified whole.

- **other**: Can be deployed independently

## 4. Roadmap

- [ ] **Monitoring**: Try connecting to Grafana Cloud and Honeycomb.io.
- [ ] **Encrypt/Decrypt**: Provides configurable encryption and decryption methods.
- [ ] **Provider**: more vault provider.
