import logging
from typing import Type

import numpy as np
from sklearn.model_selection import cross_val_score
from skopt import gp_minimize
from skopt.callbacks import DeltaYStopper
from skopt.utils import use_named_args

from lib.models import Model


def hyper_param_search(model_class: Type[Model], training):
    """
    Performs a Bayesian hyperparameter search using the given training and validation data sets.
    :param model_class: the Class object (not an instance) for the model being trained
    :param training: a pair of Numpy arrays in the format (training_data, training_labels)
    :return: scipy.optimize.OptimizeResult object containing the results of the search
    """
    logging.info("Starting hyper parameter search...")
    training_data, training_labels = training

    # Create an empty instance of the requested model class
    model = model_class.create_model()

    # Create objective function to be optimised
    @use_named_args(model_class.param_dist)
    def objective(**params):
        model.set_params(**params)
        return -np.mean(
            cross_val_score(model, training_data, training_labels, cv=2, n_jobs=1,
                            scoring='neg_mean_absolute_error')
        )

    # Run Bayesian optimization hyper parameter search using Gaussian Processes
    result = gp_minimize(objective, model_class.param_dist, callback=DeltaYStopper(5, 10), n_calls=100, verbose=True)
    return {
        'params': get_named_params(model_class.param_dist, result),
        'mean_absolute_error': result.fun
    }


def get_named_params(param_dist, result):
    logging.info("Best score: {}".format(result.fun))

    named_params = {}
    for i, param in enumerate(result.x):
        named_params[param_dist[i].name] = param

    logging.info('Best Params: {}'.format(named_params))

    return named_params
