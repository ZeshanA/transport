#!/usr/bin/env bash

export MODEL_TYPE=$1
export VENV_PATH="/data/za816/trainerEnv"
export PYTHONPATH="."

virtualenv -p /usr/bin/python3 ${VENV_PATH}
source "${VENV_PATH}/bin/activate"
pip3 install -r requirements.txt

echo "Setup complete, booting train.client to produce ${MODEL_TYPE} models..."
python3 "train/client.py" ${MODEL_TYPE}