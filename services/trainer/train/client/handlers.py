import asyncio
import logging
import multiprocessing
import sys
from concurrent.futures.thread import ThreadPoolExecutor

import gc

from lib.network import recv_json
from train.client.config import current_task_name
from train.client.requests import request_route
from train.client.tasks import tasks
from train.events import events


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
    # Send an event requesting a routeID
    await request_route(websocket)
    # Receive the response event
    await response_consumer(websocket)


async def assign_route(websocket, msg):
    loop = asyncio.get_event_loop()
    route_id = msg['routeID']
    logging.info(f"Received new Route ID: {route_id}")
    # Perform the requested task
    with ThreadPoolExecutor(max_workers=multiprocessing.cpu_count()) as executor:
        await loop.run_in_executor(executor, tasks[current_task_name], websocket, route_id)
    gc.collect()
    # Request another route
    await request_route(websocket)
    # Tell the consumer to handle the next routeID response
    await response_consumer(websocket)


async def training_complete(*_):
    logging.info("All done!")
    sys.exit(0)


# Dict containing handlers for each event type
handlers = {
    events.REGISTRATION_SUCCESS: registration_success,
    events.ASSIGN_ROUTE: assign_route,
    events.TRAINING_COMPLETE: training_complete
}
