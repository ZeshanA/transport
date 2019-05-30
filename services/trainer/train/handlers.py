import collections
import logging
import random

from lib.network import ClientSet, send_json
from lib.routes import all_routes
import train.events as events

# unprocessed_routes = collections.deque(random.sample(all_routes.copy(), len(all_routes)))
unprocessed_routes = collections.deque(all_routes[:3])


async def host_registration(websocket, client_set: ClientSet, message, path):
    host_id = message['hostID']
    logging.info(f"Received initial registration request from hostID {host_id}")
    client_set.add(host_id, websocket)
    await send_json(websocket, events.REGISTRATION_SUCCESS)


async def route_request(websocket, client_set: ClientSet, *_):
    if len(unprocessed_routes) == 0:
        await send_json(websocket, events.TRAINING_COMPLETE)
        return
    # Fetch host_id for the current socket
    host_id = client_set.get_host_id(websocket)
    # Get the next unprocessed routeID
    route_id = unprocessed_routes.pop()
    logging.info(f"Dispatching Route ID {route_id} to {host_id}")
    # Mark the routeID as assigned to the current hostID/socket
    client_set.set_route_id(socket=websocket, route_id=route_id)
    # Send an event to the client assigning it the routeID
    await send_json(websocket, events.ASSIGN_ROUTE, {"routeID": route_id})
