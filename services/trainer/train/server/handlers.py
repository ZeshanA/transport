import collections
import logging
import random

from lib.files import save_json
from lib.network import ClientSet, send_json
from lib.routes import all_routes
from train.events import events

unprocessed_routes = collections.deque(random.sample(all_routes.copy(), len(all_routes)))
model_type = None


async def host_registration(websocket, client_set: ClientSet, message, *_):
    """
    Register a host to the client set and send a registration success event.
    """
    global model_type
    host_id, model_type = message['hostID'], message['modelType']
    logging.info(f"Received initial registration request from hostID '{host_id}'")
    client_set.add(host_id, websocket)
    await send_json(websocket, events.REGISTRATION_SUCCESS)


async def route_request(websocket, client_set: ClientSet, *_):
    """
    Get the next unprocessed routeID and dispatch it to the client making the request.
    """
    if len(unprocessed_routes) == 0:
        await send_json(websocket, events.TRAINING_COMPLETE)
        return
    # Get the next unprocessed routeID
    route_id = unprocessed_routes.pop()
    # Mark the routeID as assigned to the current hostID/socket
    client_set.set_route_id(socket=websocket, route_id=route_id)
    # Send an event to the client assigning it the routeID
    await send_json(websocket, events.ASSIGN_ROUTE, {"routeID": route_id})


async def metrics_upload(websocket, client_set: ClientSet, msg, *_):
    """
    Save the performance metrics sent by the client to disk.
    """
    global model_type
    metrics, route_id = msg['metrics'], msg['metrics']['route_id']
    logging.info(f"Received performance metrics for routeID '{route_id}'")
    save_json(route_id, metrics, model_type, 'modelPerformance.json')


async def route_complete(websocket, client_set: ClientSet, *_):
    """
    Clear the completed route from being assigned to the hostID, that made the completion request.
    """
    host_id, route_id = client_set.get_host_id(websocket), client_set.get_route_id(socket=websocket)
    logging.info(f"'{host_id}' marked routeID '{route_id}' as completed")
    # Mark the routeID as completed by removing it from the host/socket
    client_set.clear_route_id(socket=websocket)
