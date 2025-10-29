# Resource Service

## 1. Overview

This project, part of **`Projectarium`**, is used to operate `resources`.

## 2. Description

There are two resource-related services: this service and **`resourcequeryservice`**. This service handles `command` in `CQRS`.

## 3. Deployment Architecture

- **Vercel**: When deployed on the `Vercel`, this service and other `Projectarium` services will be deployed as a unified whole.

- **other**: Can be deployed independently

## 4. JSON Serialization in Domain Layer

Due to Go's language constraints, implementing the `json.Marshaler` and `json.Unmarshaler` interfaces requires the methods to be defined in the same package as the type.
This means that, even though Domain-Driven Design (DDD) principles recommend keeping serialization logic in the Infrastructure layer,
in this project we place JSON serialization directly in the domain layer.

While it is technically possible to keep serialization in the Infrastructure layer using DTOs, wrappers, or mapping functions,
these approaches introduce additional complexity and boilerplate. For simplicity and practical usability,
we opted to implement serialization within the domain layer, acknowledging that this is a trade-off between strict DDD adherence and the limitations of Go.
