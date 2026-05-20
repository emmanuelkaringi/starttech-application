#!/bin/bash
set -e

echo "Deploying backend..."

# Get variables
AWS_REGION=${AWS_REGION:-us-east-1}
ECR_REPO="starttech-backend"
IMAGE_TAG=${1:-latest}

# Get AWS account ID
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_URI="$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com"

# Login to ECR
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_URI

# Build and push
cd ../backend
docker build -t $ECR_REPO:$IMAGE_TAG .
docker tag $ECR_REPO:$IMAGE_TAG $ECR_URI/$ECR_REPO:$IMAGE_TAG
docker push $ECR_URI/$ECR_REPO:$IMAGE_TAG

echo "Backend image pushed to ECR successfully!"
