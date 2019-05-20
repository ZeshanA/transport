"""
Algorithm:
    - Get path to completed models from command line
    - Loop through each folder and build up a dict of routeID -> (params, performance)
    - Calculate the weighted average using the above dict
"""
import json
import os
import sys
from glob import glob


def main():
    search_results = extract_search_results()


def extract_search_results():
    """
    Iterates over all folders (routeIDs) in the current directory and extracts
    a tuple of (bestParams, finalPerformance) from each.
    :return: a dict from routeID -> tuple(bestParams: dict, finalPerformance: dict)
    """
    path = sys.argv[1]
    os.chdir(path)
    route_ids = [p[:-1] for p in glob('*/')]
    search_results = {}
    for route_id in route_ids:
        os.chdir(route_id)
        with open('bestParams.json') as params_file:
            params = json.load(params_file)
        with open('finalPerf.json') as perf_file:
            perf = json.load(perf_file)
        search_results[route_id] = (params, perf)
        os.chdir("../")
    return search_results


if __name__ == '__main__':
    main()
