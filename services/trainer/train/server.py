import logging
import random
from typing import List, Dict

from flask import Flask, request

from lib.files import save_json
from lib.routes import all_routes

app = Flask(__name__)

unprocessed_routes: List[str] = random.sample(all_routes.copy(), len(all_routes))
currently_processing: Dict[str, str] = {}


@app.route('/getRouteID')
def get_route_id():
    host_id = request.args.get('hostID')
    if not host_id:
        logging.warning("No hostID provided by request, returning error response")
        return "Please provide a hostID in the query string"
    route_id = unprocessed_routes.pop()
    currently_processing[host_id] = route_id
    to_do, total = len(unprocessed_routes), len(all_routes)
    completed_percentage = round((total - to_do) / total, 2)
    logging.info("Assigned routeID '%s' to hostID '%s'", route_id, host_id)
    logging.info("%d of %d routes complete (%s%%)", to_do, total, completed_percentage)
    return route_id


@app.route('/completeRouteID', methods=['POST'])
def complete_route_id():
    host_id, route_id = request.args.get('hostID'), request.args.get('routeID')
    if not host_id or not route_id:
        logging.warning("Completion request is missing a hostID or a routeID, returning error response")
        return "Please provide a hostID and a routeID in the query string"
    del currently_processing[host_id]
    model_performance = request.get_json()
    save_json(route_id, model_performance, '.', 'modelPerformance.json')
    return "Marked routeID {} as complete and saved the performance info.".format(route_id)


if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s: %(levelname)s: %(message)s', level=logging.INFO)
    app.run()
