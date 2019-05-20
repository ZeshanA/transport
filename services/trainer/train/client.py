import json
import logging
import os
import random
import socket
import sys

import requests

from lib.data import get_numpy_datasets, merge_np_tuples
from lib.logs import init_logging
from lib.models import create_model, calculate_performance_metrics

SERVER_URL = "http://127.0.0.1:5000/"
GET_ROUTE_ID_URL = SERVER_URL + "getRouteID"
COMPLETE_ROUTE_ID_URL = SERVER_URL + "completeRouteID"
OPTIMAL_PARAMS = {
    'hidden_layer_count': 1,
    'neuron_count': 1,
    'activation_function': 'relu',
    'epochs': 5
}


def main():
    host_id = get_host_id()
    while True:
        # Fetch next routeID from server
        route_id = get_route_id(host_id)
        # Get train/val/test datasets
        train, val, test = get_numpy_datasets(route_id)
        # Train model using pre-determined optimal parameters
        model = get_trained_model(OPTIMAL_PARAMS, merge_np_tuples(train, val))
        # Calculate and upload final model performance metrics
        upload_performance_metrics(host_id, route_id, model, test)
        # Save model to disk
        filepath = save_model_to_disk(route_id, model)
        break


def get_host_id():
    """
    Returns a unique hostID for the current computer. This is the short hostname
    if running on the DoC network (e.g. "graphic09"), or a randomly generated
    string if not (e.g. "vast1234")
    :return: string: a unique hostID identifying the current computer to the server
    """
    hostname = parse_hostname()
    if hostname:
        return hostname
    return generate_host_id()


def parse_hostname():
    """
    Returns the shortened hostname (e.g. "graphic09") if currently executing on a DoC PC.
    None otherwise.
    :return:
    """
    hostname = socket.gethostname()
    doc_domain = ".doc.ic.ac.uk"
    if doc_domain not in hostname:
        return None
    return hostname.replace(doc_domain, '')


def generate_host_id():
    """
    Returns a unique string with the prefix "vast" to identify non-DoC hosts.
    :return: a unique hostID string to identify this PC to the server.
    """
    integer_id = str(random.randint(1, 10000))
    host_id = "vast" + integer_id
    os.environ['HOST_ID'] = host_id
    return host_id


def get_route_id(host_id):
    """
    Fetches the next routeID to process from the server, exits if the
    server says all tasks are complete.
    :param host_id: string: the hostID to report to the server
    :return: route_id: string: the next routeID to process
    """
    logging.info("Fetching next routeID from server...")
    req = requests.get(url=GET_ROUTE_ID_URL, params={'hostID': host_id})
    resp = req.text
    # Exit if there are no more routeIDs to process
    if resp == "Complete":
        logging.info("No more routeIDs to process, shutting down...")
        sys.exit()
    logging.info("Processing routeID %s...", resp)
    return resp


def get_trained_model(params, training):
    logging.info("Starting model training...")
    training_data, training_labels = training
    model = create_model(params['hidden_layer_count'], params['neuron_count'], params['activation_function'])
    model.fit(x=training_data, y=training_labels, epochs=params['epochs'])
    logging.info("Successfully completed model training...")
    return model


def upload_performance_metrics(host_id, route_id, model, test):
    metrics = calculate_performance_metrics(route_id, model, test)
    metrics_json = json.dumps(metrics)
    req = requests.post(url=COMPLETE_ROUTE_ID_URL, params={'hostID': host_id, 'routeID': route_id}, json=metrics_json)
    logging.info("Server response after performance metric submission: %s", req.text)


def save_model_to_disk(route_id, model):
    """
    Saves a trained Tensorflow model to disk. Can be loaded again using Keras'
    load_model function.
    :param route_id: string: routeID for the current model
    :param model: pointer to a trained Tensorflow model
    :return: the (relative) filepath that the model was saved at
    """
    directory = 'finalModels/{}/'.format(route_id)
    filepath = '{}/{}-finalModel.h5'.format(directory, route_id)
    os.makedirs(directory, exist_ok=True)
    model.save(filepath)
    return filepath


if __name__ == "__main__":
    init_logging()
    main()
