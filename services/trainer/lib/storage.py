import logging
import os


def get_storage_details():
    """
    Fetches the access keys for cloud object storage from the environment.
    :return: a tuple of strings (access_key_id, secret_access_key)
    """
    key_id = os.environ.get('SPACES_KEY_ID')
    secret = os.environ.get('SPACES_SECRET_KEY')
    if not key_id or not secret:
        logging.critical('SPACES_{KEY_ID/SECRET_KEY} NOT SET')
        raise KeyError
    return key_id, secret
