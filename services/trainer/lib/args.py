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
    return route_id, base_path
