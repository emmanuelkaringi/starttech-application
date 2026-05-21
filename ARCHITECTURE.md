# StartTech Application Architecture

> **🚀 Live Application**: [https://d2uv8gh2gkla6v.cloudfront.net](https://d2uv8gh2gkla6v.cloudfront.net)

Architecture documentation for the StartTech full-stack application.

---

## 📋 Table of Contents

- [System Overview](#system-overview)
- [Architecture Diagram](#architecture-diagram)
- [Component Details](#component-details)
- [Data Flow](#data-flow)
- [Authentication Flow](#authentication-flow)
- [Key Design Decisions](#key-design-decisions)
- [Security Architecture](#security-architecture)
- [Monitoring Architecture](#monitoring-architecture)

---

## System Overview

StartTech is a full-stack ToDo application built with a React frontend and Go backend API. The system is designed for high availability, scalability, and security, deployed on AWS with a complete CI/CD pipeline.

### Technology Stack

| Layer | Technology | Justification |
|-------|-----------|---------------|
| **Frontend** | React 18 + TypeScript + Vite | Type safety, fast builds, modern React features |
| **State Management** | TanStack Query + React Context | Server state caching, lightweight auth context |
| **UI Components** | Tailwind CSS + shadcn/ui | Utility-first CSS, accessible components |
| **Backend** | Go 1.25 + Gin Framework | High performance, low memory footprint |
| **Database** | MongoDB Atlas | Managed NoSQL, flexible schema for todos |
| **Caching** | Redis 7.0 (ElastiCache) | Session storage, todo list caching |
| **Authentication** | JWT with httpOnly cookies | XSS protection, stateless auth |
| **CDN** | CloudFront | Global edge caching, HTTPS termination |
| **Hosting** | S3 (Frontend) + EC2 (Backend) | Serverless static + scalable compute |
| **CI/CD** | GitHub Actions | Tight GitHub integration, free for public repos |
| **IaC** | Terraform | Declarative infrastructure, reproducible deployments |

---

## Architecture Diagram
```
┌─────────────────────────────────────────────────────────────────────────┐
│ INTERNET USERS │
└─────────────────────────────────┬───────────────────────────────────────┘
│ HTTPS
▼
┌─────────────────────────────────────────────────────────────────────────┐
│ CloudFront CDN │
│ d2uv8gh2gkla6v.cloudfront.net │
│ │
│ ┌──────────────────────┐ ┌──────────────────────────────┐ │
│ │ Static Assets │ │ API Requests │ │
│ │ (/, /.js, /.css) │ │ (/auth/, /tasks/, │ │
│ │ │ │ /users/*, /health) │ │
│ └───────────┬────────────┘ └──────────────┬───────────────┘ │
└──────────────┼────────────────────────────────────┼───────────────────────┘
│ │
▼ ▼
┌──────────────────────────┐ ┌────────────────────────────────────┐
│ S3 Bucket │ │ Application Load Balancer │
│ starttech-frontend- │ │ starttech-alb-production-xxx │
│ production-* │ │ (Port 80 → Target Group) │
└──────────────────────────┘ └────────────────┬───────────────────┘
│
▼
┌──────────────────────────────────────────────────────────────────────────┐
│ Auto Scaling Group │
│ starttech-backend-asg-production │
│ │
│ ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐ │
│ │ EC2 Instance │ │ EC2 Instance │ │ EC2 Instance │ ... │
│ │ (t3.small) │ │ (t3.small) │ │ (t3.small) │ │
│ │ │ │ │ │ │ │
│ │ ┌────────────┐ │ │ ┌────────────┐ │ │ ┌────────────┐ │ │
│ │ │ Docker │ │ │ │ Docker │ │ │ │ Docker │ │ │
│ │ │ Container │ │ │ │ Container │ │ │ │ Container │ │ │
│ │ │ :8080 │ │ │ │ :8080 │ │ │ │ :8080 │ │ │
│ │ └────────────┘ │ │ └────────────┘ │ │ └────────────┘ │ │
│ └──────────────────┘ └──────────────────┘ └──────────────────┘ │
│ │
│ Scaling Policies: │
│ - Scale Up: CPU > 80% for 10 min → +1 instance │
│ - Scale Down: CPU < 40% for 10 min → -1 instance │
│ - Min: 2 | Max: 4 | Desired: 2 │
└──────────────────────────────────────────────────────────────────────────┘
│ │ │
└────────────────────┼────────────────────┘
│
┌───────────┼──────────────────────────
│                       │            │
▼                       ▼            ▼
┌──────────────┐ ┌──────────┐ ┌──────────────┐
│ ElastiCache │ │CloudWatch│ │ MongoDB │
│ Redis 7.0 │ │Logs & │ │ Atlas │
│ (Caching) │ │Metrics │ │ (Database) │
└──────────────┘ └──────────┘ └──────────────┘
```

---

## Component Details

### Frontend (React + TypeScript)

**Location**: `frontend/`

**Key Libraries**:
| Library | Purpose |
|---------|---------|
| `react-router-dom` | Client-side routing |
| `@tanstack/react-query` | Server state management |
| `axios` | HTTP client with interceptor support |
| `tailwindcss` | Utility-first CSS framework |
| `shadcn/ui` | Accessible UI primitives |
| `zod` | Schema validation |

**State Management**:
- **Server State**: TanStack Query for API data caching, automatic refetching
- **Auth State**: React Context providing user object and auth status
- **Form State**: Local component state with zod validation

**API Client** (`src/lib/apiClient.ts`):
- Axios instance with `withCredentials: true` for cookie-based auth
- Base URL from `VITE_API_BASE_URL` environment variable
- Automatic cookie inclusion for all requests

### Backend API (Go + Gin)

**Location**: `backend/`

**Layered Architecture**:
```
cmd/api/main.go → Entry point, dependency injection
internal/routes/ → Route definitions
internal/handlers/ → HTTP request handling
internal/middleware/ → Auth, CORS, logging
internal/models/ → Data models (User, Todo)
internal/auth/ → JWT token generation/validation
internal/cache/ → Redis caching layer
internal/database/ → MongoDB connection management
internal/config/ → Environment-based configuration
internal/logger/ → Structured JSON logging
```

**Request Lifecycle**:
1. Request → CloudFront → ALB → EC2:8080
2. Gin router → CORS middleware → Logger middleware
3. Auth middleware (for protected routes) → JWT validation
4. Handler → Business logic → Database/Cache
5. Response → JSON with status code

### Database (MongoDB Atlas)

**Collections**:
| Collection | Description |
|-----------|-------------|
| `users` | User accounts (firstName, lastName, username, email, password hash) |
| `todos` | Todo items (userId, title, description, completed, timestamps) |

**Connection**: MongoDB Go Driver with connection pooling

### Caching (Redis)

**Cache Strategy**:
- **Key Pattern**: `todos:<user_id>` for user's todo lists
- **TTL**: Configurable expiration
- **Fallback**: Cache miss → Database query → Populate cache
- **Invalidation**: Cache cleared on create/update/delete operations

**Toggle**: Controlled by `ENABLE_CACHE` environment variable

---

## Data Flow

### Read Request (GET /tasks)
Browser → CloudFront → ALB → EC2 (Docker) → JWT Auth middleware

↓

(valid token)
Handler → Validate input

↓

MongoDB (insert document)

↓

Redis (invalidate cache)

↓

201 JSON Response

### Authentication Flow
**POST /auth/register**

→ Validate input
→ Check username uniqueness
→ Hash password (bcrypt)
→ Insert user into MongoDB
→ Return 201

**POST /auth/login**

→ Find user by username
→ Verify password hash
→ Generate JWT token (72h expiry)
→ Set httpOnly cookie (Secure, SameSite)
→ Return token + user info

**Subsequent Requests**

→ Auth middleware extracts JWT from cookie
→ Validates signature and expiry
→ Injects user ID into request context
→ Handler uses user ID for data access

---

## Key Design Decisions

### 1. httpOnly Cookies over Authorization Headers

**Decision**: Use httpOnly cookies instead of Bearer tokens in Authorization header.

**Rationale**:
- **XSS Protection**: JavaScript cannot access httpOnly cookies
- **Automatic Inclusion**: Browser sends cookies automatically
- **No Client-Side Token Storage**: Token never exposed to JavaScript

**Trade-off**: Requires CSRF protection (mitigated by SameSite cookie attribute)

### 2. CloudFront as API Proxy

**Decision**: Route API requests through CloudFront instead of directly to ALB.

**Rationale**:
- **Unified Domain**: Frontend and API share the same origin (no CORS preflight for same-origin)
- **HTTPS Everywhere**: CloudFront provides SSL termination
- **Edge Caching**: API responses can be cached at edge locations
- **DDoS Protection**: AWS Shield Standard included

### 3. Multi-Stage Docker Build

**Decision**: Use multi-stage Docker builds for Go application.

**Rationale**:
- **Smaller Images**: Final image ~20MB vs ~800MB
- **No Build Tools in Production**: Only the binary is copied
- **Faster Deployments**: Smaller images pull faster

### 4. Structured JSON Logging

**Decision**: Use zerolog for structured JSON logging.

**Rationale**:
- **Machine Parsable**: CloudWatch Logs Insights can query structured logs
- **Request Context**: Each log includes method, path, status, duration
- **Log Levels**: DEBUG, INFO, WARN, ERROR for filtering

---

## Security Architecture

| Layer | Implementation |
|-------|---------------|
| **Transport** | HTTPS via CloudFront (TLS 1.3) |
| **Authentication** | JWT with httpOnly, Secure, SameSite cookies |
| **Password Storage** | bcrypt hashing (cost factor 10) |
| **CORS** | Restricted to known CloudFront domain |
| **Network** | Security groups: ALB → EC2 (8080), EC2 → Redis (6379) |
| **Secrets** | GitHub Secrets for CI/CD, environment variables on EC2 |
| **IAM** | Least-privilege roles: ECR read, CloudWatch write |
| **Scanning** | Trivy vulnerability scanning in CI/CD pipeline |
| **Database** | MongoDB Atlas with TLS, IP whitelist |
| **Infrastructure** | Private subnets for EC2, public only for ALB |

---

## Monitoring Architecture

| Component | Tool | Metrics |
|-----------|------|---------|
| **Application Logs** | CloudWatch Logs | Structured JSON logs |
| **System Metrics** | CloudWatch Agent | CPU, Memory, Disk |
| **ALB Metrics** | CloudWatch | Request count, latency, errors |
| **Redis Metrics** | CloudWatch | CPU, connections, cache hit rate |
| **CloudFront** | CloudWatch | Requests, error rates |
| **Alarms** | CloudWatch Alarms | High CPU, unhealthy hosts, Redis CPU |
| **Dashboard** | CloudWatch Dashboard | Unified infrastructure view |
| **Log Analysis** | CloudWatch Logs Insights | Error analysis, performance queries |

### Log Structure
```json
{
  "level": "info",
  "method": "GET",
  "path": "/tasks",
  "status": 200,
  "duration_ms": 2.5,
  "time": "2026-05-21T12:00:00Z",
  "message": "Request completed"
}
```