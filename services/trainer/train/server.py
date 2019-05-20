import json
import logging
import random
from typing import List, Dict

from flask import Flask, request

from lib.files import save_json
from lib.logs import init_logging
from lib.routes import all_routes

app = Flask(__name__)

unprocessed_routes: List[str] = random.sample(all_routes.copy(), len(all_routes))
currently_processing: Dict[str, str] = {}


@app.route('/getRouteID')
def get_route_id():
    """
    Endpoint used by hosts to retrieve a new routeID to process.
    routeIDs are popped from the list of unprocessed routes and assigned
    to the requester's hostID in the currently_processing dict.
    :return: string: a routeID, "Complete" if all routeIDs have been processed, or an error message
    """
    host_id = request.args.get('hostID')
    if not host_id:
        logging.warning("No hostID provided by request, returning error response")
        return "Please provide a hostID in the query string"
    if not unprocessed_routes:
        return "Complete"
    route_id = unprocessed_routes.pop()
    currently_processing[host_id] = route_id
    completed, total, completed_percentage = calculate_status()
    logging.info("Assigned routeID '%s' to hostID '%s'", route_id, host_id)
    logging.info("%d of %d routes complete (%f%%)", completed, total, completed_percentage)
    return route_id


@app.route('/completeRouteID', methods=['POST'])
def complete_route_id():
    """
    Endpoint used by hosts to mark the completion of training the model for a routeID
    and to submit the final performance metrics for saving.
    :return: string: completion message or error message
    """
    host_id, route_id = request.args.get('hostID'), request.args.get('routeID')
    if not host_id or not route_id:
        logging.warning("Completion request is missing a hostID or a routeID, returning error response")
        return "Please provide a hostID and a routeID in the query string"
    del currently_processing[host_id]
    model_performance = json.loads(request.get_json())
    save_json(route_id, model_performance, '.', 'modelPerformance.json')
    return "Marked routeID {} as complete and saved the performance info.".format(route_id)


@app.route('/')
def get_status():
    """
    This endpoint prints a basic summary of current execution progress.
    :return:
    """
    completed, total, completed_percentage = calculate_status()
    response = "{} of {} routes complete ({}%)".format(completed, total, completed_percentage)
    response += ("Currently processing: " + str(currently_processing))
    return response


@app.route('/health')
def get_health():
    """
    Basic health endpoint to check if the server is alive.
    :return:
    """
    return "Healthy!"


def calculate_status():
    """
    Calculates what percentage of routes have been processed.
    :return: tuple(number of routes completed: int, total number of routes: int, percentage completed: float)
    """
    to_do, total = len(unprocessed_routes), len(all_routes)
    completed = total - to_do
    completed_percentage = round(completed / total, 3) * 100
    return completed, total, completed_percentage


if __name__ == '__main__':
    init_logging()
    app.run()
