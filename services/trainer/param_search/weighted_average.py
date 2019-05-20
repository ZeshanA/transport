import json
import os
import sys
from glob import glob

import numpy as np
from sklearn.utils.extmath import weighted_mode

TEXT_PARAMS = ['activation_function']
WEIGHT_METRIC = 'mean_absolute_error'


def main():
    params, perf = extract_search_results()
    perf_weights = get_perf_weights(perf)
    averaged_params = {}
    for param_name in params[0]:
        vals = [entry[param_name] for entry in params]
        if param_name in TEXT_PARAMS:
            averaged_params[param_name] = get_weighted_mode(vals, perf_weights)
        else:
            averaged_params[param_name] = int(np.average(vals, weights=perf_weights))
    print(averaged_params)


def extract_search_results():
    """
    Iterates over all folders (routeIDs) in the current directory and extracts
    a tuple of (bestParams, finalPerformance) from each.
    :return: a tuple of lists (bestParams, finalPerformance)
    """
    path = sys.argv[1]
    os.chdir(path)
    route_ids = [p[:-1] for p in glob('*/')]
    params_list = []
    perf_list = []
    for route_id in route_ids:
        os.chdir(route_id)
        with open('bestParams.json') as params_file:
            params = json.load(params_file)
        with open('finalPerf.json') as perf_file:
            perf = json.load(perf_file)
        params_list.append(params)
        perf_list.append(perf)
        os.chdir("../")
    return params_list, perf_list


def get_perf_weights(perf):
    perf_metrics = [entry[WEIGHT_METRIC] for entry in perf]
    max_val = max(perf_metrics)
    return list(map(lambda x: max_val - x, perf_metrics))


def get_weighted_mode(values, weights):
    return weighted_mode(values, weights)[0][0]


if __name__ == '__main__':
    main()
