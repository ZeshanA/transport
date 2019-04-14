#!/bin/sh
DEBUG=$1

if [ "${DEBUG}" = "true" ] ; then
    echo "Running Delve debugger..." && \
    /go/bin/dlv --listen=:40000 --headless=true --api-version=2 exec ./executable
else
    echo "Running executable..." && ./executable
fi
