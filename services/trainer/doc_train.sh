#!/usr/bin/env bash

export MODEL_TYPE=$2
export VENV_PATH="${HOME}/.env/transport/trainer"
export PYTHONPATH="."

source ${HOME}/.bashrc
source "${VENV_PATH}/bin/activate"
python3 --version
pip3 install -r requirements.txt

echo "Setup complete, booting train.client to produce ${MODEL_TYPE} models..."
until python3 "train/client/socket_client.py" "$@"; do
        echo "train.client crashed with exit code $?. Restarting..." >&2
        sleep 1
done