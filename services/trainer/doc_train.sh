#!/usr/bin/env bash

export MODEL_TYPE=$2
export VENV_PATH="/vol/bitbucket/za816/.env/transport/trainer"
export PYTHONPATH=".:/vol/bitbucket/za816/.local"
export HOME="/homes/za816"

# Activate new virtualenv
source "${VENV_PATH}/bin/activate" &&

# Load in environment variables from .bashrc
source "${HOME}/.bashrc" &&

# Print python3 version
python3 --version &&

# Boot crash-resilient client
echo "Setup complete, booting train.client to produce ${MODEL_TYPE} models..." &&
until python3 "train/client/socket_client.py" "$@"; do
        echo "train.client crashed with exit code $?. Restarting..." >&2
        sleep 10
done