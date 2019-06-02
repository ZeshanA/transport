import os

import boto3
import joblib

from lib.data import convert_journey_to_features
from lib.object_storage import get_storage_details


def predict(model, journey):
    feature_sample = convert_journey_to_features(journey)
    print(feature_sample)
    prediction = int(model.predict(feature_sample)[0])
    print(prediction)
    return prediction


def download_model(route_id):
    filename = 'RandomForestModel-{}-finalModel.h5'.format(route_id)
    exists = os.path.isfile(filename.format(route_id))
    if exists:
        return joblib.load(filename)
    # Model doesn't already exist on disk, download from object storage
    session = boto3.session.Session()
    key_id, secret = get_storage_details()
    client = session.client('s3',
                            region_name='fra1',
                            endpoint_url='https://fra1.digitaloceanspaces.com',
                            aws_access_key_id=key_id,
                            aws_secret_access_key=secret)
    with open(filename, 'wb') as f:
        client.download_fileobj('mtadata3', filename, f)
    return joblib.load(filename)
