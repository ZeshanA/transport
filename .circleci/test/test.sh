#!/usr/bin/env bash

# Set working directory
cd ../services

# Return non-zero exit code if any commands fail
set -e

# Extract list of services from args
SERVICES=$1

# Build and Test all services
for SERVICE in ${SERVICES}; do
    # Run `docker build` with:
    # -t (tagName): name of the service
    # -f (fileName): path to the service's Dockerfile
    docker build -t ${SERVICE} -f ${SERVICE}/Dockerfile ${SERVICE}
done