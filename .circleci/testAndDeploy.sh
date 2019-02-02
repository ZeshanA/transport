#!/usr/bin/env bash

# Set working directory
cd ../services

GO111MODULE=on

# List the directory names of your services here:
SERVICES="tsvloader"

for SERVICE in ${SERVICES}
do
    # Run `docker build` with:
    # -t (tagname): name of the service
    # -f (filename): path to the service's Dockerfile
    docker build -t ${SERVICE} -f ${SERVICE}/Dockerfile ${SERVICE}
done