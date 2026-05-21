# StartTech Application Runbook

> **🚀 Live Application**: [https://d2uv8gh2gkla6v.cloudfront.net](https://d2uv8gh2gkla6v.cloudfront.net)

Operations and troubleshooting guide for the StartTech application.

---

## 📋 Table of Contents

- [Quick Reference](#quick-reference)
- [Health Checks](#health-checks)
- [Common Operations](#common-operations)
- [Troubleshooting](#troubleshooting)
- [Emergency Procedures](#emergency-procedures)
- [Maintenance](#maintenance)

---

## Quick Reference

### Production URLs

| Service | URL |
|---------|-----|
| **Frontend** | [https://d2uv8gh2gkla6v.cloudfront.net](https://d2uv8gh2gkla6v.cloudfront.net) |
| **Backend API** | [https://d2uv8gh2gkla6v.cloudfront.net](https://d2uv8gh2gkla6v.cloudfront.net) |
| **Swagger Docs** | [https://d2uv8gh2gkla6v.cloudfront.net/swagger/index.html](https://d2uv8gh2gkla6v.cloudfront.net/swagger/index.html) |
| **Health Check** | [https://d2uv8gh2gkla6v.cloudfront.net/health](https://d2uv8gh2gkla6v.cloudfront.net/health) |
| **ALB (Direct)** | `http://starttech-alb-production-798147820.us-east-1.elb.amazonaws.com` |
| **Redis** | `starttech-redis-production.n7rpou.0001.use1.cache.amazonaws.com:6379` |

### AWS Resources

| Resource | Name |
|----------|------|
| Auto Scaling Group | `starttech-backend-asg-production` |
| Target Group | `starttech-tg-production` |
| ECR Repository | `starttech-backend` |
| S3 Bucket | `starttech-frontend-production-*` |
| CloudWatch Logs | `/starttech/application`, `/starttech/ec2/system` |

---

## Health Checks

### Quick Health Check

```bash
# Full system health
curl https://d2uv8gh2gkla6v.cloudfront.net/health
# Expected: {"cache":"ok","database":"ok"}

# Cache only: Redis connection is healthy
# Database only: MongoDB connection is healthy
# Both "ok": Application is fully operational
```

### Deep Health Check
```bash
# Check all components
./scripts/health-check.sh
```

### What Each Status Means
| Response | Meaning | Action
|----------|---------|-------
{"cache":"ok","database":"ok"} | All healthy | None
{"cache":"error","database":"ok"} | Redis down | Check ElastiCache
{"cache":"ok","database":"error"} | MongoDB down | Check Atlas status
{"cache":"error","database":"error"} | Multiple failures | Investigate immediately
No response / 502 | Backend down | See Backend 502

## Common Operations
### Deploy New Backend Version
```bash
# Option 1: Via GitHub Actions (Recommended)
# Go to Actions → Backend CI/CD Pipeline → Run workflow

# Option 2: Manual deployment
cd backend
docker build -t starttech-backend:latest .
ECR_URI="376129861708.dkr.ecr.us-east-1.amazonaws.com/starttech-backend"
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $ECR_URI
docker tag starttech-backend:latest $ECR_URI:latest
docker push $ECR_URI:latest

# Deploy to instances
./scripts/deploy-backend.sh
```

### Deploy New Frontend Version
```bash
# Option 1: Via GitHub Actions (Recommended)
# Go to Actions → Frontend CI/CD Pipeline → Run workflow

# Option 2: Manual deployment
cd frontend
npm ci
npm run build
./scripts/deploy-frontend.sh
```

### View Application Logs
```bash
# Real-time log tailing
aws logs tail /starttech/application --follow

# Last 50 lines
aws logs tail /starttech/application --follow --format short 2>&1 | head -50

# Via AWS Console
# CloudWatch → Log Groups → /starttech/application → View log streams
```

### Query Logs for Errors
```bash
# Recent errors
aws logs start-query \
  --log-group-name /starttech/application \
  --start-time $(date -d '1 hour ago' +%s) \
  --end-time $(date +%s) \
  --query-string 'fields @timestamp, @message | filter @message like /ERROR/ | sort @timestamp desc | limit 50' \
  --query 'queryId' --output text
```

### Check Running Instances
```bash
# List instances
aws autoscaling describe-auto-scaling-groups \
  --auto-scaling-group-names starttech-backend-asg-production \
  --query "AutoScalingGroups[0].Instances[?LifecycleState=='InService'].InstanceId" \
  --output text

# Check target health
TG_ARN=$(aws elbv2 describe-target-groups --names starttech-tg-production --query "TargetGroups[0].TargetGroupArn" --output text)
aws elbv2 describe-target-health --target-group-arn $TG_ARN --output table
```

### Scale the Application
```bash
# Scale up
aws autoscaling update-auto-scaling-group \
  --auto-scaling-group-name starttech-backend-asg-production \
  --desired-capacity 3

# Scale down
aws autoscaling update-auto-scaling-group \
  --auto-scaling-group-name starttech-backend-asg-production \
  --desired-capacity 2
```

### Rollback Deployment
```bash
# List recent ECR images
aws ecr describe-images --repository-name starttech-backend \
  --query "sort_by(imageDetails,&imagePushedAt)[-5:].[imagePushedAt,imageTags[0]]" \
  --output table

# Deploy a specific version
./scripts/rollback.sh
```

## Troubleshooting
### Frontend: Blank Page or 404
**Symptoms**: Page loads but is blank, or shows 404

**Causes & Fixes**:
1. SPA routing issue: CloudFront returns 403/404 for non-root paths
    - CloudFront is configured to return index.html for 403/404 errors
    - If issue persists, check CloudFront error page configuration
2. Stale cache: Old version cached
    - Invalidate CloudFront: aws cloudfront create-invalidation --distribution-id <ID> --paths "/*"
    - Wait 5-10 minutes for propagation

### Backend: 502 Bad Gateway
**Symptoms**: Frontend shows 502, health check fails

**Causes & Fixes**:

1. Container not running:

```bash
# Check container status on each instance
INSTANCE_ID="i-xxxxxxxx"
aws ssm send-command --instance-ids "$INSTANCE_ID" \
  --document-name "AWS-RunShellScript" \
  --parameters '{"commands":["docker ps -a"]}'

# Restart if needed
aws ssm send-command --instance-ids "$INSTANCE_ID" \
  --document-name "AWS-RunShellScript" \
  --parameters '{"commands":["docker start starttech-backend"]}'
```

2. Target group health check failing:

    - Container not listening on port 8080
    - Health endpoint returning non-200
    - Security group blocking ALB health checks

3. All instances unhealthy:
    - Check ECR image exists: aws ecr describe-images --repository-name starttech-backend
    - Check MongoDB connectivity
    - Redeploy via CI/CD pipeline

### Authentication: 401 Unauthorized
**Symptoms**: Login works but API calls return 401

**Causes & Fixes**:

1. **Cookie domain mismatch**: The JWT cookie domain doesn't match the frontend URL
    - Check COOKIE_DOMAINS environment variable on containers
    - Should match CloudFront domain: d2uv8gh2gkla6v.cloudfront.net
    - Cookie not being sent: Browser not including cookie
    - Clear browser cookies for the site
    - Check withCredentials: true in API client

2. **JWT token expired**: Token TTL is 72 hours by default
    - User needs to re-login
    - Check JWT_EXPIRATION_HOURS setting

### Rapid Rollback
**If a bad deployment causes issues:**

```bash
# 1. List recent ECR images
aws ecr describe-images --repository-name starttech-backend \
  --query "sort_by(imageDetails,&imagePushedAt)[-5:].[imagePushedAt,imageTags[0]]" \
  --output table

# 2. Get the previous working image tag
# 3. Update containers manually or trigger rollback script
./scripts/rollback.sh <previous-image-tag>
```

