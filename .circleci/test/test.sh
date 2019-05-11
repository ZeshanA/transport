#!/usr/bin/env bash

# Set working directory
cd ../

# Return non-zero exit code if any commands fail
set -e

# Extract list of services from args by trimming off the first arg (which is the program name)
SERVICES=${@:1}

# Build and Test all services
for SERVICE in ${SERVICES}; do
    # Run `docker build` with:
    # -t (tagName): name of the service
    # -f (fileName): path to the service's Dockerfile
    echo "TESTING ${SERVICE} FROM [${SERVICES}]"
    docker build -t ${SERVICE} -f services/${SERVICE}/Dockerfile --build-arg SERVICE_NAME=${SERVICE} .
done

# Test lib folder
cd lib && go test ./...
cd ../