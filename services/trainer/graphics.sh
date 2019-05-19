#!/usr/bin/env bash

export ROUTE_ID=$1
export VENV_PATH="/data/za816/trainerEnv"
echo "Processing ${ROUTE_ID}..."

virtualenv -p /usr/bin/python3 ${VENV_PATH}
source "${VENV_PATH}/bin/activate"
pip3 install -r requirements.txt

echo "Setup complete, starting param search..."
python3 "param_search/search.py" "${ROUTE_ID}" graphics