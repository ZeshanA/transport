import logging


def init_logging():
    logging.basicConfig(format='%(asctime)s: %(levelname)s: %(message)s')
    logging.getLogger().setLevel(logging.INFO)
