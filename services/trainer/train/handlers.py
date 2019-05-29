from lib.network import ClientSet, send_json
from train.events import REGISTRATION_SUCCESS


async def host_registration(websocket, client_set: ClientSet, message, path):
    host_id = message['hostID']
    print(f"Received initial registration request from hostID {host_id}")
    websocket.host_id = host_id
    client_set.add(host_id, websocket)
    await send_json(websocket, REGISTRATION_SUCCESS)

