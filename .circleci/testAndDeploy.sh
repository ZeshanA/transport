#!/usr/bin/env bash

# Set working directory
cd ../services

GO111MODULE=on

# List the directory names of your services here:
SERVICES="tsvloader"

# Build and Test all services
for SERVICE in ${SERVICES}; do
    # Run `docker build` with:
    # -t (tagName): name of the service
    # -f (fileName): path to the service's Dockerfile
    docker build -t ${SERVICE} -f ${SERVICE}/Dockerfile ${SERVICE}
done

# Deploy only services that have changed

# Extract commit range (or single commit hash)
COMMIT_RANGE=$(echo "${CIRCLE_COMPARE_URL}" | cut -d/ -f7)

# If there was only a single commit, convert it to range format (i.e. add `...`)
if [[ ${COMMIT_RANGE} != *"..."* ]]; then
  COMMIT_RANGE="${COMMIT_RANGE}...${COMMIT_RANGE}"
fi

AZURE_USERNAME="za816"
AZURE_RESOURCE_GROUP="transport"
AZURE_PLAN="transportPlan"

IMAGE_REPO_USERNAME="zshnamjd"
CHANGED_FILES=$(git diff --name-only ${COMMIT_RANGE} | tr "\n" " ")
echo "Changed files: $CHANGED_FILES"

if [[ ${CIRCLE_BRANCH} != "nycData" ]]; then
    echo "Branch ${CIRCLE_BRANCH} does not need to be deployed"
else
    echo "Deploying changed services to production"
    for SERVICE in ${SERVICES}; do
        SHOULD_DEPLOY=false

        for CHANGED_FILE in ${CHANGED_FILES}; do
            if [[ ${CHANGED_FILE} = *${SERVICE}* ]]; then
                echo "The service '${SERVICE}' has changed - will deploy"
                SHOULD_DEPLOY=true
                break
            fi
        done

        if ${SHOULD_DEPLOY}; then
            echo "Deploying ${SERVICE}"
            docker build -t ${SERVICE} -f ${SERVICE}/Dockerfile ${SERVICE}
            docker tag ${SERVICE}:latest ${IMAGE_REPO_USERNAME}/${SERVICE}:latest
            docker push ${IMAGE_REPO_USERNAME}/${SERVICE}:latest
            az webapp deployment user set --user-name ${AZURE_USERNAME} --password ${AZURE_PASSWORD}
            az webapp create --resource-group ${AZURE_RESOURCE_GROUP} --plan ${AZURE_RESOURCE_GROUP} --name ${SERVICE} \
            --deployment-container-image-name ${IMAGE_REPO_USERNAME}/${SERVICE} && echo "Deployed ${SERVICE} successfully"
        else
            echo "No changes for service '${SERVICE}' - skipping deployment"
        fi
    done
    echo "Deployment complete"
fi