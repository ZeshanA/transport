#!/usr/bin/env bash

export VENV_PATH="/data/za816/trainerEnv"
export PYTHONPATH="."

virtualenv -p /usr/bin/python3 ${VENV_PATH}
source "${VENV_PATH}/bin/activate"
pip3 install -r requirements.txt

echo "Setup complete, starting train.server..."
python3 "train/server.py" & (sleep 10 && echo "Starting train.client..." && python3 "train/client.py")