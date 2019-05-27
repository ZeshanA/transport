import logging
import os
import sys

import boto3

sys.path.insert(0, os.path.abspath('.'))

from lib.data import get_numpy_datasets, merge_np_tuples
from lib.logs import init_logging
from lib.models import create_model, calculate_performance_metrics

OPTIMAL_PARAMS = {
    'hidden_layer_count': 47,
    'neuron_count': 552,
    'activation_function': 'relu',
    'epochs': 30
}


def main():
    # Fetch next routeID from server
    route_id = sys.argv[1]
    # Get train/val/test datasets
    train, val, test = get_numpy_datasets(route_id)
    # Train model using pre-determined optimal parameters
    model = get_trained_model(OPTIMAL_PARAMS, merge_np_tuples(train, val))
    # Calculate and upload final model performance metrics
    metrics = calculate_performance_metrics(route_id, model, test)
    # Save model to disk
    filepath = save_model_to_disk(route_id, model)
    # Upload model to object storage
    upload_model(route_id, filepath)
    print(metrics)


def save_model_to_disk(route_id, model):
    """
    Saves a trained Tensorflow model to disk. Can be loaded again using Keras'
    load_model function.
    :param route_id: string: routeID for the current model
    :param model: pointer to a trained Tensorflow model
    :return: the (relative) filepath that the model was saved at
    """
    directory = '/data/za816/trained/models/{}/'.format(route_id)
    filepath = '{}/finalModel.h5'.format(directory, route_id)
    os.makedirs(directory, exist_ok=True)
    model.save(filepath)
    return filepath


def get_trained_model(params, training):
    """
    Creates a model based on the given parameters, trains it on the provided training data
    and returns a pointer to the fully trained model.
    :param params: a dict containing the parameters that the model should be trained using
    :param training: a tuple of Numpy arrays (training_data, training_labels)
    :return: a pointer to a trained Tensorflow model
    """
    logging.info("Starting model training...")
    training_data, training_labels = training
    model = create_model(params['hidden_layer_count'], params['neuron_count'], params['activation_function'])
    model.fit(x=training_data, y=training_labels, epochs=params['epochs'])
    logging.info("Successfully completed model training...")
    return model


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


def upload_model(route_id, filepath):
    """
    Uploads the .h5 file stored at the given filepath to cloud object storage.
    :param route_id: string: the routeID that the model corresponds to
    :param filepath: string: the path to the .h5 file containing the trained Tensorflow model
    """
    logging.info("Uploading final model for routeID %s to storage...", route_id)
    session = boto3.session.Session()
    key_id, secret = get_storage_details()
    client = session.client('s3',
                            region_name='fra1',
                            endpoint_url='https://fra1.digitaloceanspaces.com',
                            aws_access_key_id=key_id,
                            aws_secret_access_key=secret)
    client.upload_file(filepath, 'mtadata2', '{}-finalModel.h5'.format(route_id), ExtraArgs={'ACL': 'public-read'})
    logging.info("Successfully uploaded final model for routeID %s to storage...", route_id)


def print_separator():
    """
    Prints an ASCII divider to separate different executions in the console
    """
    print("\n\n===================================================================================================\n\n")


if __name__ == "__main__":
    init_logging()
    main()
