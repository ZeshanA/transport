from lib.hosts import get_host_id
from lib.network import send_json
from train.client.config import get_model_type
from train.events import events


async def register_host_id(websocket):
    """
    Sends a hostID registration event on the provided websocket.
    """
    host_id = get_host_id()
    await send_json(websocket, events.START_REGISTRATION, {"hostID": host_id, "modelType": get_model_type()})


async def request_route(websocket):
    """
    Sends a routeID request event on the provided websocket.
    """
    await send_json(websocket, events.ROUTE_REQUEST)


async def upload_performance_metrics(websocket, metrics):
    """
    Uploads performance metrics for the current model to the websocket.
    """
    await send_json(websocket, events.METRICS_UPLOAD, {'metrics': metrics})


async def upload_best_parameter_set(websocket, result):
    """
    Uploads performance metrics for the current model to the websocket.
    """
    await send_json(websocket, events.PARAMETER_SET_UPLOAD, result)


async def complete_route(websocket, route_id):
    """
    Sends a route completion request to the websocket provided.
    """
    await send_json(websocket, events.ROUTE_COMPLETE, {'route_id': route_id})
