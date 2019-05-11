from typing import List

import logging

import routes
import sys

from db.db import connect


def main():
    logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(levelname)s: %(message)s')
    execute_mode()


def execute_mode():
    mode: str = sys.argv[1]
    if mode == "single":
        route: str = sys.argv[2]
        logging.info("Training only routeID {}".format(route))
        train(route)
    elif mode == "all":
        logging.info("Training all routes sequentially")
        train_all(routes.all_routes)
    else:
        logging.error("{} is not a valid mode, choose either 'single' or 'all'".format(mode))
        exit()


def train(route: str):
    logging.info("Training on routeID: {}".format(route))
    conn = connect()
    # TODO: Pull in all data for this routeID from labelled_journeys
    conn.close()


def train_all(route_ids: List[str]):
    for route in route_ids:
        train(route)


if __name__ == "__main__":
    main()
