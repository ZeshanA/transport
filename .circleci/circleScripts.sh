#!/usr/bin/env bash

# Set working directory
cd .circleci

# Return non-zero exit code if any commands fail
set -e

# List the directory names of your services here:
SERVICES_TO_TEST="tsvloader"
SERVICES_TO_DEPLOY=""

# Run test or deploy script based on command arg
if [[ $1 == "test" ]]; then
    ./test/test.sh ${SERVICES_TO_TEST}
elif [[ $1 == "deploy" ]]; then
    ./deploy/deploy.sh ${SERVICES_TO_DEPLOY}
else
    echo "Unrecognised command '${1}': you can run 'test' or 'deploy'"
fi
