#!/bin/bash

# Build and push Battlesnake to Azure Container Registry
# Usage: ./build_and_push.sh <registry-name> [image-tag]

set -e

# Configuration
REGISTRY_NAME="${1:-}"
IMAGE_TAG="${2:-latest}"

# Validate inputs
if [ -z "$REGISTRY_NAME" ]; then
    echo "Usage: ./build_and_push.sh <registry-name> [image-tag]"
    echo ""
    echo "Example: ./build_and_push.sh erikregistry latest"
    echo "Example: ./build_and_push.sh erikregistry v1.0.0"
    exit 1
fi

# Remove .azurecr.io suffix if provided
REGISTRY_NAME="${REGISTRY_NAME%.azurecr.io}"
# Remove https:// prefix if provided
REGISTRY_NAME="${REGISTRY_NAME#https://}"

# Set variables
REGISTRY_URL="${REGISTRY_NAME}.azurecr.io"
IMAGE_NAME="battlesnake"
FULL_IMAGE_NAME="${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG}"

echo "Building Battlesnake Docker image..."
echo "Registry: $REGISTRY_URL"
echo "Image: $FULL_IMAGE_NAME"
echo ""

# Build the Docker image
docker build -t "$FULL_IMAGE_NAME" .

if [ $? -ne 0 ]; then
    echo "Docker build failed!"
    exit 1
fi

echo ""
echo "Logging into Azure Container Registry..."

# Login to Azure Container Registry
az acr login --name "$REGISTRY_NAME"

if [ $? -ne 0 ]; then
    echo "Failed to login to ACR. Make sure you're authenticated with 'az login'"
    exit 1
fi

echo ""
echo "Pushing image to Azure Container Registry..."

# Push to ACR
docker push "$FULL_IMAGE_NAME"

if [ $? -ne 0 ]; then
    echo "Docker push failed!"
    exit 1
fi

echo ""
echo "✓ Successfully pushed $FULL_IMAGE_NAME"
echo ""
echo "Image details:"
echo "  Registry: $REGISTRY_URL"
echo "  Image: $IMAGE_NAME"
echo "  Tag: $IMAGE_TAG"
