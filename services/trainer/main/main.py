import logging
import sys
from typing import List

from db.db import connect


def main():
    logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(levelname)s: %(message)s')
    route_id = extract_route_id()
    train(route_id)


def extract_route_id():
    if len(sys.argv) < 2:
        logging.critical("No route ID was provided, correct usage: ./trainer <routeID>")
        sys.exit(1)
    return sys.argv[1]


def train(route: str):
    logging.info("Training on routeID: {}".format(route))
    conn = connect()
    rows = get_rows_for_route_id(route, conn)
    print(rows)
    conn.close()


def get_rows_for_route_id(route_id: str, db_conn):
    logging.info("Fetching rows for routeID: {}".format(route_id))
    # Execute query to fetch all rows for the given route ID
    cur = db_conn.cursor()
    # TODO: Remove LIMIT 10
    query = "SELECT * FROM labelled_journey2 WHERE line_ref=%s LIMIT 10;"
    cur.execute(query, (route_id,))
    rows = cur.fetchall()
    # Close cursor and return fetched rows
    cur.close()
    logging.info("{} rows successfully fetched".format(len(rows)))
    return rows


def train_all(route_ids: List[str]):
    for route in route_ids:
        train(route)


if __name__ == "__main__":
    main()
