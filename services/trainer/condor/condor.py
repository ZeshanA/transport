#!/usr/bin/env python3
import logging
import sys

import subprocess

from lib.routes import all_routes


def main():
    host_id, host_count, base_path = int(sys.argv[1]), int(sys.argv[2]), sys.argv[3] if len(sys.argv) > 3 else '.'
    task_count = len(all_routes) // host_count
    first_task_index = host_id * task_count
    for i in range(first_task_index, first_task_index + task_count):
        route_id = all_routes[i]
        subprocess.call(["python", "param_search/search.py", route_id, base_path])


if __name__ == "__main__":
    logging.basicConfig()
    logging.getLogger().setLevel(logging.INFO)
    main()
