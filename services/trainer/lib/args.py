import logging
import sys


# Extracts the route_id from the CLI arguments, exiting if none was provided
def extract_route_id():
    if len(sys.argv) < 2:
        logging.critical("No route ID was provided, correct usage: ./trainer <routeID>")
        sys.exit(1)
    return sys.argv[1]
