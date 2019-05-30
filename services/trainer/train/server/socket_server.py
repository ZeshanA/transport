import asyncio
import json
import logging

import websockets

from lib.logs import init_logging
from lib.network import ClientSet
from train.events import events
from train.server.handlers import host_registration, route_request, metrics_upload, unprocessed_routes, route_complete

# Set of currently connected clients, accessible by hostID,
# their live websocket object, or the routeID they're currently assigned to.
connected_clients = ClientSet()

# Dict of async handler functions for each event type
handlers = {
    events.START_REGISTRATION: host_registration,
    events.ROUTE_REQUEST: route_request,
    events.METRICS_UPLOAD: metrics_upload,
    events.ROUTE_COMPLETE: route_complete
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
    except websockets.ConnectionClosed:
        # Closed/failed connections are okay: the finally block will perform clean up and ensure
        # any incomplete routes are still trained, no need to raise an exception here.
        pass
    finally:
        # Client has disconnected (gracefully or abruptly)
        host_id = connected_clients.get_host_id(websocket)
        route_id = connected_clients.get_route_id(host_id)
        # If the client hadn't finished processing its routeID, add the routeID back to the queue
        if route_id is not None:
            logging.warning(f"{host_id} did not finish training for {route_id}, adding it back to the pool")
            unprocessed_routes.append(route_id)
        # Remove the client from the connected set
        connected_clients.remove(host_id=host_id)
        logging.info(f"Connection closed by hostID '{host_id}'")


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
