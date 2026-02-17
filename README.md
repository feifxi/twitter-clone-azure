# Twitter Clone Azure

This is a full-stack Twitter clone application built with a modern tech stack, designed for scalability and deployment on Microsoft Azure.

## Overview

The project consists of two main components:
- [**Frontend**](./twitter-frontend/README.md): A Next.js 15 application with a responsive UI.
- [**Backend**](./twitter-backend/README.md): A Spring Boot REST API backed by PostgreSQL.

## Features

- **User Authentication**: Secure sign-up/login.
- **Tweet Management**: Create, view, like, retweet, reply, and delete tweets.
- **Media Support**: Image uploads for tweets and profile pictures.
- **Real-time Updates**: Optimistic UI updates for interactions.
- **Search**: Users and hashtags.
- **Responsive Design**: Mobile-friendly interface.

## Quick Start (Local Development)

### Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose installed.

## Quick Start (Local Development)

### Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose installed.
- Node.js 18+ and Java 17+ (if running services outside Docker).

### Running Locally

1.  **Start Database**:
    The `docker-compose.yml` is configured to run only the PostgreSQL database for local development.
    ```bash
    docker-compose up -d
    ```
    This starts a PostgreSQL instance on port `5432`.

2.  **Start Backend**:
    Navigate to `twitter-backend` and run:
    ```bash
    ./mvnw spring-boot:run
    ```

3.  **Start Frontend**:
    Navigate to `twitter-frontend` and run:
    ```bash
    npm run dev
    ```

## Deployment (Azure & GitHub Actions)

This project uses **Azure Container Apps** for hosting and **GitHub Actions** for CI/CD.

### Architecture
- **Frontend**: Next.js App -> Azure Container App
- **Backend**: Spring Boot App -> Azure Container App
- **Database**: Azure Database for PostgreSQL (Flexible Server)
- **Modularity**: Frontend and Backend have separate CI/CD pipelines triggered only by changes in their respective directories.

### Setup Guide

1.  **Azure Resources**:
    - Create an Azure Container Registry (ACR).
    - Create an Azure Database for PostgreSQL (Flexible Server).
    - Create an Azure Container Apps Environment.

2.  **GitHub Secrets**:
    Add the following secrets to your GitHub repository:
    - `AZURE_CREDENTIALS`: JSON output from `az ad sp create-for-rbac`.
    - `ACR_LOGIN_SERVER`: e.g., `myregistry.azurecr.io`.
    - `ACR_USERNAME` & `ACR_PASSWORD`: Access keys for your ACR.
    - `AZURE_RESOURCE_GROUP`: Your Resource Group name.
    - `NEXT_PUBLIC_API_URL`: URL of your deployed Backend Container App.
    - `SPRING_DATASOURCE_URL`, `USERNAME`, `PASSWORD`: Your Azure Postgres credentials.
    - `JWT_SECRET`: A strong secret for token generation.

3.  **Deploy**:
    Push changes to the `main` branch.
    - Changes in `twitter-frontend/**` trigger the Frontend Pipeline.
    - Changes in `twitter-backend/**` trigger the Backend Pipeline.

## License

MIT
