# Twitter Clone Monorepo

Full-stack Twitter/X clone with a Go-first backend and Next.js frontend.

The Go API is the primary backend in active development:
- `twitter-go-api` (primary)
- `twitter-next-web` (frontend)
- `twitter-java-api` (legacy/alternative backend)

## Stack

- Frontend: Next.js (App Router), TypeScript, Tailwind, TanStack Query, Zustand
- Primary Backend: Go, Gin, PostgreSQL + sqlc, Redis, Azure Blob Storage
- Legacy Backend: Java (Spring Boot), PostgreSQL

## Local Development

### 1. Start infra (Postgres + Redis)

```bash
docker-compose up -d
```

### 2. Run Go API

```bash
cd twitter-go-api
make migrateup
make run
```

### 3. Run Next.js Web

```bash
cd twitter-next-web
npm install
npm run dev
```

Web: `http://localhost:3000`  
API: `http://localhost:8080`

## Validation Commands

Go API:
```bash
cd twitter-go-api
go test ./...
```

Next.js Web:
```bash
cd twitter-next-web
npx tsc --noEmit
npm run lint
```

## Deployment

Deployment is configured with GitHub Actions to Azure Container Apps.

Workflows:
- `.github/workflows/deploy-go-api.yml`
  - Trigger: push to `main` when `twitter-go-api/**` changes, or manual dispatch
  - Deploy target: `ca-chntwt-go-api-dev`
- `.github/workflows/deploy-next-web.yml`
  - Trigger: push to `main` when `twitter-next-web/**` changes, or manual dispatch
  - Deploy target: `ca-chntwt-web-dev`
- `.github/workflows/deploy-java-api.yml`
  - Trigger: manual dispatch only
  - Deploy target: `ca-chntwt-api-dev`

All pipelines use:
- `azure/login@v2` (OIDC)
- `azure/container-apps-deploy-action@v2`
- Azure Container Registry: `crchntwtdevsea001.azurecr.io`
