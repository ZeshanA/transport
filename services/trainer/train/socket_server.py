import asyncio
import json

import websockets

from lib.logs import init_logging
from lib.network import ClientSet
from train.events import START_REGISTRATION
from train.handlers import host_registration

connected_clients = ClientSet()
handlers = {
    START_REGISTRATION: host_registration
}


def main():
    # Start the websocket server and ensure it runs continuously
    start_server = websockets.serve(consumer_handler, '0.0.0.0', 8765)
    asyncio.get_event_loop().run_until_complete(start_server)
    asyncio.get_event_loop().run_forever()


async def consumer_handler(websocket, path):
    """
    First handler called when a websocket connection is made.
    :param websocket: the websocket object for the new connection
    :param path: the path that the connection was made from
    :return: void
    """
    try:
        # Send messages to the consumer as they come in
        async for message in websocket:
            await consumer(websocket, message, path)
    finally:
        print(f"Connection closed for Host ID: {connected_clients.get_host_id(websocket)}")


async def consumer(websocket, message, path):
    # Decode JSON message into dict
    json_msg = json.loads(message)
    # Extract the event that occurred
    event = json_msg['event']
    # Call the correct handler for the event
    await handlers[event](websocket, connected_clients, json_msg, path)


if __name__ == '__main__':
    init_logging()
    main()
