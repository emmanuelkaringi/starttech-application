# StartTech Application

Full-stack application with React frontend and Go backend.

## Project Structure
```
starttech-application/
├── .github/workflows/
│ ├── backend-ci-cd.yml # Backend CI/CD pipeline
│ └── frontend-ci-cd.yml # Frontend CI/CD pipeline
├── frontend/ # React application (Client)
├── backend/ # Go API server (Server/MuchToDo)
└── scripts/
├── deploy-frontend.sh
├── deploy-backend.sh
├── health-check.sh
└── rollback.sh
```

## CI/CD Pipelines

### Backend Pipeline
- Runs unit tests and code quality checks
- Builds Docker image with security scanning
- Pushes to Amazon ECR
- Deploys to EC2 Auto Scaling Group

### Frontend Pipeline
- Installs dependencies and runs tests
- Builds production React bundle
- Runs security audit
- Deploys to S3 bucket
- Invalidates CloudFront cache

## Environment Variables

### Backend
Required environment variables:
- `MONGO_URI` - MongoDB Atlas connection string
- `JWT_SECRET_KEY` - Secret key for JWT tokens
- `REDIS_ADDR` - Redis endpoint address

### Frontend
- `VITE_API_URL` - Backend API URL

## Deployment

### Prerequisites
- AWS CLI configured
- ECR repository created
- Infrastructure deployed via starttech-infra

### Manual Deployment
```bash
# Backend
cd backend
docker build -t starttech-backend .
# Push to ECR and deploy

# Frontend
cd frontend
npm run build
# Sync dist/ to S3 bucket
