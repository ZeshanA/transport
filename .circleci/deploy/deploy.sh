#!/usr/bin/env bash

# Set working directory
cd deploy

# Return non-zero exit code if any commands fail
set -e

# Extract list of services from args
SERVICES=$1

# Extract commit range (or single commit hash)
LAST_SUCCESSFUL_BUILD_URL="https://circleci.com/api/v1.1/project/github/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}/tree/${CIRCLE_BRANCH}?filter=completed&limit=1"
LAST_SUCCESSFUL_COMMIT=$(curl -Ss -u "${CIRCLE_API_KEY}:" ${LAST_SUCCESSFUL_BUILD_URL} | jq -r '.[0]["vcs_revision"]')
CHANGED_FILES=$(git diff --name-only ${CIRCLE_SHA1}..${LAST_SUCCESSFUL_COMMIT} | tr "\n" " ")

echo "Changed files: $CHANGED_FILES"

# Only deploy master branch
if [[ ${CIRCLE_BRANCH} != "master" ]]; then
    echo "Branch ${CIRCLE_BRANCH} does not need to be deployed"
else
    # Install azure-cli
    chmod +x ./installAzureCLI.sh && ./installAzureCLI.sh

    # Only deploy services where changes were detected
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
            chmod +x ./deploySingleService.sh && ./deploySingleService.sh ${SERVICE}
        else
            echo "No changes for service '${SERVICE}' - skipping deployment"
        fi
    done
    echo "Deployment complete"
fi