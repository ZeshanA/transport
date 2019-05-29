import asyncio
import json
import time

import websockets

from train.client import get_host_id
from train.events import START_REGISTRATION

WEBSOCKET_SERVER = 'ws://localhost:8765'


async def hello():
    async with websockets.connect(WEBSOCKET_SERVER) as websocket:
        host_id = get_host_id()
        msg = json.dumps({
            "event": START_REGISTRATION,
            "hostID": host_id
        })
        await websocket.send(msg)

        response = await websocket.recv()
        print(f"Server response: {response}")


asyncio.get_event_loop().run_until_complete(hello())
