import logging
import os
import random
import socket
import sys

import requests

sys.path.insert(0, os.path.abspath('.'))

from models.neural_network import NNModel
from models.random_forest import RandomForestModel
from lib.data import get_numpy_datasets
from lib.logs import init_logging, print_separator

# TODO: Change server URL
# SERVER_URL = "http://d.zeshan.me:5000/"
SERVER_URL = "http://127.0.0.1:5000/"
GET_ROUTE_ID_URL = SERVER_URL + "getRouteID"
COMPLETE_ROUTE_ID_URL = SERVER_URL + "completeRouteID"

MODEL_TYPES = {
    'neural_network': NNModel,
    'random_forest': RandomForestModel
}


def main():
    host_id, model_type_id = get_host_id(), get_model_type()
    logging.info("The hostID for this machine is %s\n", host_id)
    model_class = MODEL_TYPES[model_type_id]
    while True:
        # Fetch next routeID from server
        route_id = get_route_id(host_id)
        # Get train/val/test datasets
        train, test = get_numpy_datasets(route_id, False)
        # Create the requested model
        model = model_class(route_id)
        # Train the model
        model.train(train)
        # Calculate and upload final model performance metrics
        model.upload_performance_metrics(COMPLETE_ROUTE_ID_URL, host_id, test)
        # Upload model to object storage
        model.upload_model()
        # Print ASCII divider for clarity in console
        print_separator()


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
    Returns the previously set host_id if it exists.
    None otherwise.
    :return:
    """
    if os.getenv('HOST_ID') is not None:
        return os.environ['HOST_ID']
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


def get_model_type():
    """
    :return: a string identifying the type of model to be trained, e.g. "neural_network"
    """
    if len(sys.argv) < 2:
        raise ValueError(
            """
            Please pass in the type of model you wish to train as a command line parameter.
            You can pick from the following values: {}
            """.format(list(MODEL_TYPES.keys()))
        )
    return sys.argv[1]


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


if __name__ == "__main__":
    init_logging()
    main()
