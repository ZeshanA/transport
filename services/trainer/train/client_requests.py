from lib.network import send_json
from train import events
from train.client import get_host_id


async def register_host_id(websocket):
    """
    Sends a hostID registration event on the provided websocket
    """
    host_id = get_host_id()
    await send_json(websocket, events.START_REGISTRATION, {"hostID": host_id})


async def request_route(websocket):
    await send_json(websocket, events.ROUTE_REQUEST)
