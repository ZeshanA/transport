import json
import logging

from flask import Flask, request

import stop_to_stop
from lib.logs import init_logging
from lib.model import download_model, predict

app = Flask(__name__)


@app.route('/predictFromMovement', methods=['POST'])
def get_prediction_from_movement():
    journey = request.get_json()
    route_id = journey['LineRef']
    logging.info(f"Received movement prediction request for routeID {route_id}")
    model = download_model(route_id)
    return json.dumps({
        "prediction": predict(model, journey)
    })


@app.route('/predictStopToStop', methods=['POST'])
def get_stop_to_stop_prediction():
    req = request.get_json()
    mvmt = req['sampleMovement']
    route_id = mvmt['LineRef']
    logging.info(f"Received stop-to-stop prediction request for routeID {route_id}")
    model = download_model(route_id)
    prediction = stop_to_stop.calculate(model, req)
    return json.dumps({
        "prediction": prediction
    })


if __name__ == '__main__':
    init_logging()
    app.run(host='0.0.0.0')
