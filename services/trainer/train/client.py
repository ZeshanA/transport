import logging
import os
import random
import socket
import sys

import requests

from lib.data import get_numpy_datasets, merge_np_tuples
from lib.logs import init_logging
from lib.models import create_model, save_performance_metrics, calculate_performance_metrics

SERVER_URL = "http://127.0.0.1:5000/"
GET_ROUTE_ID_URL = SERVER_URL + "getRouteID"
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
        metrics = calculate_performance_metrics(route_id, model, test)
        break


def get_trained_model(params, training):
    logging.info("Starting model training...")
    training_data, training_labels = training
    model = create_model(params['hidden_layer_count'], params['neuron_count'], params['activation_function'])
    model.fit(x=training_data, y=training_labels, epochs=params['epochs'])
    logging.info("Successfully completed model training...")
    return model


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


if __name__ == "__main__":
    init_logging()
    main()
