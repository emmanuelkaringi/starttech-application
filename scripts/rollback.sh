#!/bin/bash
set -e

echo "Rolling back..."

# Get previous Docker image tag from ECR
AWS_REGION=${AWS_REGION:-us-east-1}
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_URI="$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com"

# List recent images and get the previous one
echo "Available images in ECR:"
aws ecr describe-images --repository-name starttech-backend --query "sort_by(imageDetails,&imagePushedAt)[-3:].[imagePushedAt,imageTags[0]]" --output table

echo "To rollback, run the deploy script with a specific image tag:"
echo "./deploy-backend.sh <image-tag>"
