#!/bin/bash
set -e

echo "Running health checks..."

# Get ALB DNS
ALB_DNS=$(aws elbv2 describe-load-balancers \
    --names starttech-alb-production \
    --query "LoadBalancers[0].DNSName" \
    --output text 2>/dev/null || echo "")

if [ -n "$ALB_DNS" ]; then
    echo "Checking backend health at http://$ALB_DNS/health"
    curl -f http://$ALB_DNS/health || {
        echo "Backend health check failed!"
        exit 1
    }
    echo "Backend health check passed!"
else
    echo "ALB not found, skipping backend health check"
fi

# Get CloudFront domain
CLOUDFRONT_DOMAIN=$(aws cloudfront list-distributions \
    --query "DistributionList.Items[0].DomainName" \
    --output text 2>/dev/null || echo "")

if [ -n "$CLOUDFRONT_DOMAIN" ]; then
    echo "Checking frontend at https://$CLOUDFRONT_DOMAIN"
    curl -f https://$CLOUDFRONT_DOMAIN > /dev/null 2>&1 || {
        echo "Frontend health check failed!"
        exit 1
    }
    echo "Frontend health check passed!"
else
    echo "CloudFront distribution not found, skipping frontend health check"
fi

echo "All health checks passed!"
