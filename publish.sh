#!/bin/bash

DOCKER_USERNAME="niklas1simakov"
IMAGE_VERSION="1.0.0"

# Log in to Docker Hub
echo "Logging in to Docker Hub..."
docker login

# Define services and platforms
services=("frontend" "get-books" "post-books" "put-books" "delete-books")
PLATFORMS="linux/amd64,linux/arm64"

# Build and push multi-platform images
for service in "${services[@]}"; do
  echo "Building and pushing multi-platform image for $DOCKER_USERNAME/cc-ex-3-${service}:$IMAGE_VERSION and :latest..."
  docker buildx build \
    --platform $PLATFORMS \
    --tag $DOCKER_USERNAME/cc-ex-3-${service}:$IMAGE_VERSION \
    --tag $DOCKER_USERNAME/cc-ex-3-${service}:latest \
    --file Dockerfile \
    --target ${service}-service \
    --push .

  if [ $? -ne 0 ]; then
    echo "Error building/pushing ${service}. Exiting."
    exit 1
  fi
done

echo "All multi-platform images built, tagged, and pushed successfully!"

echo "\n--- List of pushed image manifests (check Docker Hub for full details) ---"
for service in "${services[@]}"; do
  echo "Manifest for $DOCKER_USERNAME/cc-ex-3-${service}:$IMAGE_VERSION:"
  docker buildx imagetools inspect $DOCKER_USERNAME/cc-ex-3-${service}:$IMAGE_VERSION || true
  echo "Manifest for $DOCKER_USERNAME/cc-ex-3-${service}:latest:"
  docker buildx imagetools inspect $DOCKER_USERNAME/cc-ex-3-${service}:latest || true
done 