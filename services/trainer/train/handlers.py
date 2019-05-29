import collections
import logging
import random

from lib.network import ClientSet, send_json
from lib.routes import all_routes
import train.events as events

unprocessed_routes = collections.deque(random.sample(all_routes.copy(), len(all_routes)))


async def host_registration(websocket, client_set: ClientSet, message, path):
    host_id = message['hostID']
    logging.info(f"Received initial registration request from hostID {host_id}")
    client_set.add(host_id, websocket)
    await send_json(websocket, events.REGISTRATION_SUCCESS)


async def route_request(websocket, client_set: ClientSet, *_):
    if len(unprocessed_routes) == 0:
        await send_json(websocket, events.TRAINING_COMPLETE)
        return
    route_id = unprocessed_routes.pop()
    logging.info(f"Dispatching Route ID {route_id} to {client_set.get_host_id(websocket)}")
    await send_json(websocket, events.ASSIGN_ROUTE, {"routeID": route_id})
