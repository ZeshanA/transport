import os
import random
import socket
import sys

import requests

SERVER_URL = "http://127.0.0.1:5000/"
GET_ROUTE_ID_URL = SERVER_URL + "getRouteID"


def main():
    host_id = get_host_id()
    route_id = get_route_id(host_id)


def get_route_id(host_id):
    """
    Fetches the next routeID to process from the server, exits if the
    server says all tasks are complete.
    :param host_id: string: the hostID to report to the server
    :return: route_id: string: the next routeID to process
    """
    req = requests.get(url=GET_ROUTE_ID_URL, params={'hostID': host_id})
    resp = req.text
    # Exit if there are no more routeIDs to process
    if resp == "Complete":
        sys.exit()
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
    main()
