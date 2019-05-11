from functools import singledispatch
from typing import List

import logging
import routes
import sys


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
        train(routes.all_routes)
    else:
        logging.error("{} is not a valid mode, choose either 'single' or 'all'".format(mode))
        exit()


@singledispatch
# TODO: Implement
def train(route_ids: List[str]):
    print("List of route IDs")


@train.register
# TODO: Implement
def _(route: str):
    print("Single route ID")


if __name__ == "__main__":
    main()
