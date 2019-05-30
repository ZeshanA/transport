import logging
import sys


def extract_route_id() -> str:
    """
    Extracts the route_id from the CLI arguments, exiting if none was provided
    :return: route_id (string)
    """
    if len(sys.argv) < 2:
        logging.critical("No route ID was provided, correct usage: ./trainer <routeID>")
        sys.exit(1)
    return sys.argv[1]


def extract_cli_args():
    """
    Extract route_id and base_path arguments from sys.argv
    :return: route_id, base_path
    """
    route_id = extract_route_id()
    base_path = sys.argv[2] if len(sys.argv) > 2 else '.'
    logging.info("route_id = {}, base_path = {}".format(route_id, base_path))
    return route_id, base_path


def get_model_type():
    """
    :return: a string identifying the type of model to be trained, e.g. "neural_network"
    """
    if len(sys.argv) < 2:
        raise ValueError('Please pass in the type of model you wish to train as a command line parameter.')
    return sys.argv[1]
