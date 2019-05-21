import json
import logging
import os
from datetime import datetime

import boto3 as boto3
import pandas as pd
import sklearn
from flask import Flask, request
from pytz import timezone
from sklearn_pandas import gen_features, DataFrameMapper
from tensorflow.keras.models import load_model

app = Flask(__name__)

NUMERIC_COLS = ["direction_ref", "longitude", "latitude", "distance_from_stop",
                "day", "month", "year", "hour", "minute", "second", "estimate"]
TEXT_COLS = ["operator_ref", "progress_rate", "occupancy", "stop_point_ref"]
COL_COUNT = len(NUMERIC_COLS) + len(TEXT_COLS)


@app.route('/', methods=['POST'])
def get_prediction():
    journey = request.get_json()
    route_id = journey['LineRef']
    model = download_model(route_id)
    feature_sample = convert_journey_to_features(journey)
    prediction = int(model.predict(feature_sample)[0][0])
    return json.dumps({
        "prediction": prediction
    })


def get_storage_details():
    """
    Fetches the access keys for cloud object storage from the environment.
    :return: a tuple of strings (access_key_id, secret_access_key)
    """
    key_id = os.environ.get('SPACES_KEY_ID')
    secret = os.environ.get('SPACES_SECRET_KEY')
    if not key_id or not secret:
        logging.critical('SPACES_{KEY_ID/SECRET_KEY} NOT SET')
        raise KeyError
    return key_id, secret


def download_model(route_id):
    session = boto3.session.Session()
    key_id, secret = get_storage_details()
    client = session.client('s3',
                            region_name='fra1',
                            endpoint_url='https://fra1.digitaloceanspaces.com',
                            aws_access_key_id=key_id,
                            aws_secret_access_key=secret)
    filename = '{}-finalModel.h5'.format(route_id)
    exists = os.path.isfile(filename.format(route_id))
    if exists:
        return load_model(filename)
    with open(filename, 'wb') as f:
        client.download_fileobj('mtadata', filename, f)
    return load_model(filename)


def convert_journey_to_features(journey):
    df = journey_to_dataframe(journey)
    feature_def = gen_features(
        columns=TEXT_COLS,
        classes=[sklearn.preprocessing.LabelEncoder]
    )
    mapper = DataFrameMapper(feature_def, default=None)
    sample = mapper.fit_transform(df)
    return sample


def journey_to_dataframe(journey):
    timestamp = parse_datetime(journey['Timestamp'])
    eat = journey['ExpectedArrivalTime']
    if eat is None:
        estimate = None
    else:
        eat = parse_datetime(eat)
        estimate = (eat - timestamp).total_seconds()
    extracted_fields = {
        "direction_ref": [int(journey['DirectionRef'])],
        "operator_ref": [journey['OperatorRef']],
        "longitude": [float(journey['Longitude'])],
        "latitude": [float(journey['Latitude'])],
        "progress_rate": [journey['ProgressRate']],
        "occupancy": [journey['Occupancy']],
        "distance_from_stop": [int(journey['DistanceFromStop'])],
        "stop_point_ref": [journey['StopPointRef']],
        "day": [int(timestamp.strftime('%d'))],
        "month": [int(timestamp.strftime('%m'))],
        "year": [int(timestamp.strftime('%Y'))],
        "hour": [int(timestamp.strftime('%H'))],
        "minute": [int(timestamp.strftime('%M'))],
        "second": [int(timestamp.strftime('%S'))],
        "estimate": [estimate]
    }
    return pd.DataFrame.from_dict(extracted_fields)


def parse_datetime(string):
    date_format, est = '%Y-%m-%d %H:%M:%S', timezone('US/Eastern')
    naive = datetime.strptime(string, date_format)
    return est.localize(naive)


if __name__ == '__main__':
    app.run()
