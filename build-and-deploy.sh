#!/bin/bash

# Variables
AWS_REGION="us-east-1" # e.g. us-east-1
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
FRONTEND_REPO_NAME="minecraft-control-panel-frontend"
BACKEND_REPO_NAME="minecraft-control-panel-backend"
FRONTEND_PATH="./frontend-control" # Replace with the actual path to your frontend Dockerfile
BACKEND_PATH="./backend"   # Replace with the actual path to your backend Dockerfile
FRONTEND_IMAGE_TAG="latest"
BACKEND_IMAGE_TAG="latest"

docker buildx create --use || true

# Create the ECR repositories
echo "Creating ECR repositories..."

# aws ecr create-repository --repository-name $FRONTEND_REPO_NAME --region $AWS_REGION
# aws ecr create-repository --repository-name $BACKEND_REPO_NAME --region $AWS_REGION

# Get the ECR login
echo "Logging in to Amazon ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

# Build the Docker images
echo "Building Docker images..."

# Frontend image
docker buildx build --platform linux/amd64 -t ${FRONTEND_REPO_NAME}:${FRONTEND_IMAGE_TAG} $FRONTEND_PATH --load

# Backend image
docker buildx build --platform linux/amd64 -t ${BACKEND_REPO_NAME}:${BACKEND_IMAGE_TAG} $BACKEND_PATH --load

# Tag the images for ECR
echo "Tagging Docker images for ECR..."

# Tag frontend
docker tag ${FRONTEND_REPO_NAME}:${FRONTEND_IMAGE_TAG} ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${FRONTEND_REPO_NAME}:${FRONTEND_IMAGE_TAG}

# Tag backend
docker tag ${BACKEND_REPO_NAME}:${BACKEND_IMAGE_TAG} ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${BACKEND_REPO_NAME}:${BACKEND_IMAGE_TAG}

# Push the images to ECR
echo "Pushing Docker images to ECR..."

# Push frontend
docker push ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${FRONTEND_REPO_NAME}:${FRONTEND_IMAGE_TAG}

# Push backend
docker push ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${BACKEND_REPO_NAME}:${BACKEND_IMAGE_TAG}

echo "Docker images successfully pushed to ECR!"