import logging
import random
from typing import List, Dict

from flask import Flask, request

from lib.routes import all_routes

app = Flask(__name__)

unprocessed_routes: List[str] = random.sample(all_routes.copy(), len(all_routes))
currently_processing: Dict[str, str] = {}


@app.route('/getRouteID')
def get_route_id():
    host_id = request.args.get('hostID')
    if not host_id:
        logging.warning("No host ID provided by request, returning error response")
        return "Please provide a host ID in the query string"
    route_id = unprocessed_routes.pop()
    currently_processing[host_id] = route_id
    to_do, total = len(unprocessed_routes), len(all_routes)
    completed_percentage = round((total - to_do) / total, 2)
    logging.info("Assigned routeID '%s' to hostID '%s'", route_id, host_id)
    logging.info("%d of %d routes complete (%s%%)", to_do, total, completed_percentage)
    return route_id


if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s: %(levelname)s: %(message)s', level=logging.INFO)
    app.run()
