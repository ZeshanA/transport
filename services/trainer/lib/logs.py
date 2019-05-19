import logging


def init_logging():
    for handler in logging.root.handlers[:]:
        logging.root.removeHandler(handler)
    logging.basicConfig(
        format='%(asctime)s [%(filename)s/%(funcName)s:%(lineno)d] %(levelname)s – %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S',
        level=logging.INFO
    )
