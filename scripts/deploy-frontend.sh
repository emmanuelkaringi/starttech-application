#!/bin/bash
set -e

echo "Deploying frontend..."

# Build frontend
cd ../frontend
npm ci
npm run build

# Get bucket name
BUCKET_NAME=$(aws s3api list-buckets --query "Buckets[?starts_with(Name, 'starttech-frontend-production')].Name" --output text)

if [ -z "$BUCKET_NAME" ]; then
    echo "Error: S3 bucket not found"
    exit 1
fi

# Sync to S3
aws s3 sync dist/ s3://$BUCKET_NAME/ --delete

echo "Frontend deployed to S3 bucket: $BUCKET_NAME"
