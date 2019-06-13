import logging


def init_logging():
    """
    Initialises the root logger to use the following convenient custom format:
     "2019-05-29 21:05:10 [logs.py/init_logging:7] INFO – Your message of choice."
    """
    for handler in logging.root.handlers[:]:
        logging.root.removeHandler(handler)
    logging.basicConfig(
        format='%(asctime)s [%(filename)s/%(funcName)s:%(lineno)d] %(levelname)s – %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S',
        level=logging.INFO
    )