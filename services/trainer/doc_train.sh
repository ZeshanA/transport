#!/usr/bin/env bash

export MODEL_TYPE=$2
export VENV_PATH="/data/za816/.env/transport/trainer"
export PYTHONPATH="."

# Copy pyenv installation over to /data/ to avoid quota
if [[ ! -d "/data/za816/.pyenv" ]]; then
    echo "Copying pyenv install to data directory..."
    mkdir -p "/data/za816/.pyenv" && cp -R -n "${HOME}/.pyenv" "/data/za816"
    echo "Copying pyenv complete..."
fi

# Create new virtualenv using the interpreter copied over
if [[ ! -d "${VENV_PATH}" ]]; then
    virtualenv -p "/data/za816/.pyenv/versions/3.7.3/bin/python3" ${VENV_PATH}
fi

# Activate new virtualenv
source "${VENV_PATH}/bin/activate"

# Load in environment variables from .bashrc
source "${HOME}/.bashrc"

# Print python3 version
python3 --version

# Install dependencies inside virtualenv
pip3 install -r requirements.txt

# Boot crash-resilient client
echo "Setup complete, booting train.client to produce ${MODEL_TYPE} models..."
until python3 "train/client/socket_client.py" "$@"; do
        echo "train.client crashed with exit code $?. Restarting..." >&2
        sleep 10
done