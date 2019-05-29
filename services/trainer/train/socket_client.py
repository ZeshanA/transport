import asyncio
import logging
import sys
import time

import websockets

import train.events as events
from lib.logs import init_logging
from lib.network import send_json, recv_json
from train.client import get_host_id

WEBSOCKET_SERVER = 'ws://localhost:8765'


def main():
    asyncio.get_event_loop().run_until_complete(run_client())


async def run_client():
    async with websockets.connect(WEBSOCKET_SERVER) as websocket:
        await register_host_id(websocket)
        registration_msg = await recv_json(websocket)
        # TODO: Switch to dict-based dispatch to a handler or a switch statement, this is getting messy already
        if registration_msg['event'] == events.REGISTRATION_SUCCESS:
            while True:
                await request_route(websocket)
                resp = await recv_json(websocket)
                if resp['event'] == events.ASSIGN_ROUTE:
                    logging.info(f"Received new Route ID: {resp['routeID']}")
                    # TODO: Perform the training task in a new thread
                    time.sleep(0.5)
                elif resp['event'] == events.TRAINING_COMPLETE:
                    logging.info("All done!")
                    sys.exit(0)
        elif registration_msg['event'] == events.TRAINING_COMPLETE:
            logging.info("All done!")
            sys.exit(0)


async def register_host_id(websocket):
    host_id = get_host_id()
    await send_json(websocket, events.START_REGISTRATION, {"hostID": host_id})


async def request_route(websocket):
    await send_json(websocket, events.ROUTE_REQUEST)


if __name__ == '__main__':
    init_logging()
    main()
