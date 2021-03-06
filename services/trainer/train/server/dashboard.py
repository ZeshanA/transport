import threading

from flask import Flask, render_template

from lib.network import ClientSet
from lib.routes import all_routes
from train.server.handlers import unprocessed_routes
import train.server.config as config

app = Flask(__name__)

clients: ClientSet = ClientSet()


@app.route('/')
def get_status():
    """
    This endpoint prints a basic summary of current execution progress.
    :return:
    """
    completed, total, completed_percentage = calculate_status()
    # response += ("Currently processing: " + clients.current_state_string())
    return render_template('dashboard.html',
                           model_type=format_model_type(config.model_type),
                           completed=completed, total=total, completed_percentage=completed_percentage,
                           assigned=trim_host_ids(clients.current_state()))


def format_model_type(mt):
    """
    Converts the model type identifier into a human-readable string.
    e.g. "random_forest" -> "random forest"
    :param mt: string: the model type ID
    :return: a human-readable version with underscores stripped and capitalised words
    """
    components = mt.split('_')
    formatted_components = list(map(lambda word: word.capitalize(), components))
    return ' '.join(formatted_components)


def calculate_status():
    """
    Calculates what percentage of routes have been processed.
    :return: tuple(number of routes completed: int, total number of routes: int, percentage completed: float)
    """
    connected_count = clients.connected_hosts_count()
    unassigned, total = len(unprocessed_routes), len(all_routes)
    completed = total - unassigned - connected_count
    completed_percentage = '{number:.{digits}f}'.format(number=(completed / total) * 100, digits=2)
    return completed, total, completed_percentage


def trim_host_ids(current_state):
    result, max_length = {}, 12
    for host_id, route_id in current_state.items():
        trimmed_host_id = host_id[:max_length]
        result[trimmed_host_id] = route_id
    return result


def start(connected_clients: ClientSet):
    global clients
    clients = connected_clients
    # Spin up a new thread to run the Flask server
    t = threading.Thread(target=app.run, kwargs={'host': '0.0.0.0', 'port': 5000})
    t.start()
