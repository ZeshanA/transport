#!/usr/bin/env bash

source ~/.env/transport/trainer/bin/activate

export hostID=$1
export hostCount=$2

python condor.py ${hostID} ${hostCount}