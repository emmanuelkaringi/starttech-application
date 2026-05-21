# StartTech Application Runbook

## Quick Reference

| Component | URL |
|-----------|-----|
| Production Frontend | https://d37rt95zzt3h33.cloudfront.net |
| Backend API (via CF) | https://d37rt95zzt3h33.cloudfront.net |
| Backend API (direct) | http://starttech-alb-production-1469746205.us-east-1.elb.amazonaws.com |
| Swagger Docs | https://d37rt95zzt3h33.cloudfront.net/swagger/index.html |
| Health Check | https://d37rt95zzt3h33.cloudfront.net/health |

## Common Commands

### Local Development

```bash
# Start backend locally
cd backend
docker-compose up -d mongodb redis  # Start dependencies
go run ./cmd/api/main.go            # Start API server

# Start frontend locally
cd frontend
npm run dev                         # Start Vite dev server
```

## Build Docker Image
```sh
cd backend
docker build -t starttech-backend:latest .
```

## Run Tests
```sh
# Backend unit tests
cd backend && go test ./... -v -short

# Backend integration tests
INTEGRATION=true go test -v --tags=integration ./...

# Frontend tests
cd frontend && npm test
```

## Environment Variables
### Backend Required Variables

Variable	Description	Example
MONGO_URI	MongoDB connection string	mongodb+srv://...
JWT_SECRET_KEY	JWT signing secret	Random 32-byte hex
REDIS_ADDR	Redis endpoint	host:6379
ENABLE_CACHE	Enable Redis cache	true/false
ALLOWED_ORIGINS	CORS allowed origins	https://domain.com
SECURE_COOKIE	Use Secure cookies	true for HTTPS
COOKIE_DOMAINS	Cookie domain	cloudfront.net

### Frontend Required Variables
Variable	Description
VITE_API_BASE_URL	Backend API URL