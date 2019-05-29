import json
from typing import Dict


class ClientSet:
    def __init__(self):
        self.sockets_by_host_id = {}
        self.route_ids_by_host_id = {}
        self.host_ids_by_socket = {}

    def add(self, host_id, socket):
        self.sockets_by_host_id[host_id], self.host_ids_by_socket[socket] = socket, host_id

    def remove(self, host_id=None, socket=None):
        if socket:
            host_id = self.host_ids_by_socket[socket]
        del self.host_ids_by_socket[socket]
        del self.route_ids_by_host_id[host_id]
        del self.sockets_by_host_id[host_id]

    def get_socket(self, host_id):
        return self.sockets_by_host_id[host_id]

    def get_host_id(self, socket):
        return self.host_ids_by_socket[socket]

    def get_route_id(self, host_id=None, socket=None):
        if socket:
            host_id = self.host_ids_by_socket[socket]
        return self.route_ids_by_host_id[host_id]

    def set_route_id(self, host_id=None, socket=None, route_id=None):
        if socket:
            host_id = self.host_ids_by_socket[socket]
        self.route_ids_by_host_id[host_id] = route_id


async def send_json(websocket, event: str, msg: Dict = None):
    """
    Send a JSON event with optional additional fields via the given websocket connection.
    :param websocket: the websocket to send the message on
    :param event: the desired value of the "event" field inside the JSON message
    :param msg: a dict containing any additional fields for the JSON message to contain
    :return:
    """
    if msg is None:
        msg = {}
    msg['event'] = event
    json_msg = json.dumps(msg)
    await websocket.send(json_msg)


async def recv_json(websocket):
    response = await websocket.recv()
    return json.loads(response)
