import asyncio

import websockets

from lib.logs import init_logging
from train.client.handlers import response_consumer
from train.client.requests import register_host_id

WEBSOCKET_SERVER = 'ws://localhost:8765'


async def run_client():
    """
    Creates a websocket connection to the server, registers our hostID
    and continues processing the next event in the response received.
    """
    async with websockets.connect(WEBSOCKET_SERVER) as websocket:
        await register_host_id(websocket)
        await response_consumer(websocket)


def main():
    """
    Run client until it completes.
    """
    asyncio.get_event_loop().run_until_complete(run_client())


if __name__ == '__main__':
    init_logging()
    main()
