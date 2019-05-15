import json
import os


def save_json(route_id, filename, object):
    """
    Saves the given object to models/{routeID}/{filename} in JSON format, creating intermediary folders
    and overwriting the existing file if necessary.
    :param route_id: the route id currently being calculated
    :param filename: the name of the file to save the object under
    :param object: a JSON-serialisable object to be saved at models/{route_id}/filename
    :return: void
    """
    dir = "models/{}/".format(route_id)
    filepath = dir + filename
    # Create directory if needed
    if not os.path.exists(dir):
        os.makedirs(dir)
    # Delete existing file if needed
    if os.path.exists(filepath):
        os.remove(filepath)
    # Write best params to file in JSON format
    file = open(filepath, 'w+')
    file.write(json.dumps(object))
    file.close()
