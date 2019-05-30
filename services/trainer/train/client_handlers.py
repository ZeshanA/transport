import logging
import sys
import time

from lib.network import recv_json
from train import events
from train.client_requests import request_route


async def response_consumer(websocket):
    """
    Waits for a response from the websocket, converts it into
    JSON and dispatches it to the correct handler based on the event
    type.
    :param websocket:
    """
    json_msg = await recv_json(websocket)
    await handlers[json_msg['event']](websocket, json_msg)


async def registration_success(websocket, msg):
    while True:
        # Send an event requesting a routeID
        await request_route(websocket)
        # Receive the response event
        await response_consumer(websocket)


async def assign_route(websocket, msg):
    logging.info(f"Received new Route ID: {msg['routeID']}")
    # TODO: Perform the training task in a new thread
    time.sleep(0.5)
    await request_route(websocket)
    await response_consumer(websocket)


async def training_complete(*_):
    logging.info("All done!")
    sys.exit(0)


# Dict of handlers for each event type
handlers = {
    events.REGISTRATION_SUCCESS: registration_success,
    events.ASSIGN_ROUTE: assign_route,
    events.TRAINING_COMPLETE: training_complete
}
