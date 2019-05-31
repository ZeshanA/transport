import logging
import sys


# Config Helpers

def get_task_name():
    if len(sys.argv) < arg_count:
        fatal_arg_error()
    return sys.argv[task_arg_index]


def get_model_type():
    """
    :return: a string identifying the type of model to be trained, e.g. "neural_network"
    """
    if len(sys.argv) < arg_count:
        fatal_arg_error()
    return sys.argv[model_arg_index]


def fatal_arg_error(msg="Please provide all args, usage: python3 train/socket_client.py <task> <model_type>"):
    logging.info(msg)
    sys.exit(1)


# Config Globals

arg_count = 3
task_arg_index, model_arg_index = 1, 2

current_task_name = get_task_name()
