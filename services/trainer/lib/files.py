import json
import os


def save_json(route_id, object, base_path, filename):
    """
    Saves the given object to models/{routeID}/{filename} in JSON format, creating intermediary folders
    and overwriting the existing file if necessary.
    :param route_id: the route id currently being calculated
    :param object: a JSON-serialisable object to be saved at models/{route_id}/filename
    :param base_path: a string path to the parent folder of the models directory
    :param filename: the name of the file to save the object under
    :return: void
    """
    directory = "{}/models/{}/".format(base_path, route_id)
    filepath = directory + filename
    # Create directory if needed
    if not os.path.exists(directory):
        os.makedirs(directory)
    # Delete existing file if needed
    if os.path.exists(filepath):
        os.remove(filepath)
    # Write best params to file in JSON format
    file = open(filepath, 'w+')
    file.write(json.dumps(object))
    file.close()
