# StartTech Application

> **Live Application**: [https://d2uv8gh2gkla6v.cloudfront.net](https://d2uv8gh2gkla6v.cloudfront.net)

Full-stack ToDo application with React frontend, Go backend API, Redis caching, and MongoDB database. Complete CI/CD pipeline with GitHub Actions deploying to AWS.

---

## Table of Contents

- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Quick Links](#quick-links)
- [Prerequisites](#prerequisites)
- [Local Development](#local-development)
- [Environment Variables](#environment-variables)
- [CI/CD Pipelines](#cicd-pipelines)
- [Deployment](#deployment)
- [Testing](#testing)
- [API Documentation](#api-documentation)
- [Infrastructure](#infrastructure)

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Frontend** | React 18, TypeScript, Vite, Tailwind CSS, shadcn/ui |
| **Backend** | Go 1.21, Gin Framework |
| **Database** | MongoDB Atlas |
| **Caching** | Redis (ElastiCache) |
| **Auth** | JWT with httpOnly cookies |
| **CI/CD** | GitHub Actions |
| **Hosting** | S3 + CloudFront (Frontend), EC2 + ALB (Backend) |
| **Monitoring** | CloudWatch Logs, Metrics, Alarms |
| **Infrastructure** | Terraform (see [starttech-infra](https://github.com/Innocent9712/starttech-infra)) |

---

## Project Structure
```
starttech-application/
├── .github/workflows/
│ ├── frontend-ci-cd.yml # Frontend CI/CD pipeline
│ └── backend-ci-cd.yml # Backend CI/CD pipeline
├── frontend/ # React + TypeScript (Vite)
│ ├── src/
│ │ ├── components/ # UI components (shadcn/ui)
│ │ ├── context/ # React Context (Auth)
│ │ ├── hooks/ # Custom hooks
│ │ ├── lib/ # API client, utilities
│ │ ├── routes/ # Page components
│ │ └── types/ # TypeScript types
│ ├── .env.example # Environment template
│ └── package.json
├── backend/ # Go API (Gin)
│ ├── cmd/api/ # Application entry point
│ ├── internal/
│ │ ├── auth/ # JWT token service
│ │ ├── cache/ # Redis cache layer
│ │ ├── config/ # Configuration loader
│ │ ├── database/ # MongoDB connection
│ │ ├── handlers/ # HTTP handlers
│ │ ├── logger/ # Structured logging
│ │ ├── middleware/ # Auth, CORS, logging
│ │ ├── models/ # Data models
│ │ └── routes/ # Route definitions
│ ├── Dockerfile # Multi-stage Docker build
│ └── .env.example # Environment template
├── scripts/
│ ├── deploy-frontend.sh # Frontend deployment helper
│ ├── deploy-backend.sh # Backend deployment helper
│ ├── health-check.sh # Health check script
│ └── rollback.sh # Rollback helper
├── ARCHITECTURE.md # System architecture
├── RUNBOOK.md # Operations guide
└── README.md # This file
```

---

## Quick Links

| Service | URL |
|---------|-----|
| **Frontend** | [https://d2uv8gh2gkla6v.cloudfront.net](https://d2uv8gh2gkla6v.cloudfront.net) |
| **Backend API** | [https://d2uv8gh2gkla6v.cloudfront.net](https://d2uv8gh2gkla6v.cloudfront.net) (via CloudFront) |
| **Swagger Docs** | [https://d2uv8gh2gkla6v.cloudfront.net/swagger/index.html](https://d2uv8gh2gkla6v.cloudfront.net/swagger/index.html) |
| **Health Check** | [https://d2uv8gh2gkla6v.cloudfront.net/health](https://d2uv8gh2gkla6v.cloudfront.net/health) |

---

## 📋 Prerequisites

### For Local Development
- **Node.js** >= 18
- **Go** >= 1.25
- **Docker** & Docker Compose
- **MongoDB** (Atlas)
- **AWS CLI** (for deployment)

### For CI/CD
- GitHub repository secrets configured
- AWS account with appropriate permissions
- ECR repository created

---

## Local Development

### Backend Setup

```bash
cd backend

# Copy environment template
cp .env.example .env
# Edit .env with your local values

# Start dependencies (MongoDB + Redis)
docker-compose up -d

# Install Go dependencies
go mod tidy

# Generate API documentation
swag init -g cmd/api/main.go

# Run the server (http://localhost:8080)
go run ./cmd/api/main.go
```

### Frontend Setup

```bash
cd frontend

# Copy environment template
cp .env.example .env

# Install dependencies
npm install

# Start dev server (http://localhost:5173)
npm run dev
```

## Environment Variables
### Backend


| Variable   | Required | Default | Description
|------------|----------|---------|------------|
| PORT       | No	| 8080	| Server port
| MONGO_URI	| Yes	| -	| MongoDB Atlas connection string
| DB_NAME	| No	| | much_todo_db	Database name
| JWT_SECRET_KEY	| Yes	| -	| JWT signing secret (use openssl rand -hex 32)
| JWT_EXPIRATION_HOURS	| No	| 72 |Token expiration time
| ENABLE_CACHE	| No	| false | Enable Redis caching
| REDIS_ADDR	| No	| localhost:6379 | Redis endpoint
| ALLOWED_ORIGINS	| No	| http://localhost:5173 | CORS allowed origins
| SECURE_COOKIE	| No	| false | Set Secure flag on cookies
| COOKIE_DOMAINS	| No	| - | Allowed cookie domains
| LOG_LEVEL	| No	| INFO | Logging level (DEBUG/INFO/WARN/ERROR)
| LOG_FORMAT	| No	| json | Log format (json/text)

### Frontend

| Variable   | Required | Description
|------------|----------|---------|
VITE_API_BASE_URL       |Yes | Backend API URL

## CI/CD Pipelines
### Backend Pipeline (`.github/workflows/backend-ci-cd.yml`)

Triggered on push to `main`, changes to `backend/**`

| Stage | Description
|-------|------------|
Test | Go unit tests, code quality checks, vulnerability scan
Build & Push | Docker build, image scan, push to ECR
Deploy | Rolling update to EC2 via SSM, smoke tests

### Frontend Pipeline (`.github/workflows/frontend-ci-cd.yml`)

Triggered on push to main, changes to frontend/**

| Stage | Description
|-------|------------|
Build & Test | npm install, lint, test, security audit, production build
Deploy | Sync to S3, CloudFront cache invalidation

## Deployment
### Automatic (CI/CD)
Push to `main` branch triggers the appropriate pipeline automatically.

### Manual Deployment
```bash
# Backend
cd backend
docker build -t starttech-backend:latest .
# Tag and push to ECR, then deploy via scripts/deploy-backend.sh

# Frontend
cd frontend
npm run build
# Sync dist/ to S3 bucket via scripts/deploy-frontend.sh
```

## Testing
### Backend

```bash
cd backend

# Unit tests
go test ./... -v -short

# Integration tests (requires Docker)
INTEGRATION=true go test -v --tags=integration ./...
```

### Frontend
```bash
cd frontend

# Run tests
npm test

# Run linting
npm run lint
```

## API Documentation
**Interactive Swagger documentation is available at:**

**Production**: https://d2uv8gh2gkla6v.cloudfront.net/swagger/index.html

**Local**: http://localhost:8080/swagger/index.html

### API Endpoints

|Method | Path | Description | Auth
|-------|------|-------------|-----|
GET | /health | Health check | No
POST | /auth/register | Register user | No
POST | /auth/login | Login | No
POST | /auth/logout | Logout | No
GET | /tasks | Get all todos | Yes
POST | /tasks | Create todo | Yes
GET | /tasks/:id | Get todo by ID | Yes
PUT | /tasks/:id | Update todo | Yes
DELETE | /tasks/:id | Delete todo | Yes
GET | /users/me | Get profile | Yes
PUT | /users/me | Update profile | Yes
DELETE | /users/me | Delete account | Yes

## Infrastructure

**Infrastructure is managed in a separate repository:**

[starttech-infra]((https://github.com/emmanuelkaringi/starttech-infra)) - Terraform configurations, monitoring setup, CI/CD for infrastructure

### Key infrastructure components:

1. VPC with public/private subnets across 2 AZs

2. EC2 Auto Scaling Group (2-4 instances)

3. Application Load Balancer

4. CloudFront CDN with S3 and ALB origins

5. ElastiCache Redis 7.0

6. CloudWatch Logs, Metrics, Alarms

7. MongoDB Atlas for database

## Security
1. **Secrets**: All credentials stored in GitHub Secrets, never in code

2. **Authentication**: JWT with httpOnly, Secure cookies

3. **HTTPS**: Enforced via CloudFront

4. **CORS**: Restricted to known origins

5. **IAM**: Least-privilege roles for EC2 instances

6. **Scanning**: Trivy vulnerability scanning in CI/CD

7. **Network**: Security groups restricting traffic between components