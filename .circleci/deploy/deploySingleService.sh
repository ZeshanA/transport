#!/usr/bin/env bash

# Set authentication vars
AZURE_USERNAME="za816"
AZURE_RESOURCE_GROUP="transport"
AZURE_PLAN="transportPlan"

# Extract the service to deploy from args
SERVICE=$1

# Return non-zero exit code if any commands fail
set -e

echo "Deploying ${SERVICE}"

# Set working directory
cd ../../services

# Login to Docker Hub
echo ${DOCKER_PASSWORD} | docker login -u ${DOCKER_USERNAME} --password-stdin

# Tag the service's local image with the repository path needed to push it to Docker Hub
docker tag ${SERVICE}:latest ${DOCKER_USERNAME}/${SERVICE}:latest

# Push image to Docker Hub
docker push ${DOCKER_USERNAME}/${SERVICE}:latest

# Login to Azure and web app deployment service
az login --username ${AZURE_USERNAME}@ic.ac.uk --password ${AZURE_PASSWORD}
az webapp deployment user set --user-name ${AZURE_USERNAME} --password ${AZURE_PASSWORD}

# Deploy the Docker Hub image that we pushed earlier to Docker
az webapp create --resource-group ${AZURE_RESOURCE_GROUP} --plan ${AZURE_PLAN} --name ${SERVICE} \
--deployment-container-image-name ${DOCKER_USERNAME}/${SERVICE} && echo "Deployed ${SERVICE} successfully"

# Restart the container
az webapp restart --resource-group ${AZURE_RESOURCE_GROUP} --name ${SERVICE}