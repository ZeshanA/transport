import logging
from typing import Type

from keras.wrappers.scikit_learn import KerasRegressor
from sklearn.ensemble import RandomForestRegressor
from sklearn.model_selection import RandomizedSearchCV

from lib.models import Model


def hyper_param_search(model_class: Type[Model], training, validation):
    """
    Performs a randomised grid search using the given training and validation data sets.
    :param model_class: the Class object (not an instance) for the model being trained
    :param training: a pair of Numpy arrays in the format (training_data, training_labels)
    :param validation: a pair of Numpy arrays in the format (validation_data, validation_labels)
    :return: SciKit.cv_results_ object containing the results of the search
    """
    logging.info("Starting hyper parameter search...")
    training_data, training_labels = training

    # Define the type of model we'll be using
    model = KerasRegressor(build_fn=model_class.create_model)

    # Get iterable ranges for each of our hyperparameters
    param_dist = model_class.param_dist

    if not param_dist:
        return None

    # Define the parameters for the search itself
    iteration_count = 4
    random_search = RandomizedSearchCV(RandomForestRegressor(),
                                       param_distributions=param_dist,
                                       n_iter=iteration_count,
                                       n_jobs=-1,
                                       cv=2,
                                       verbose=3)

    # Perform the search and return the results
    result = random_search.fit(training_data, training_labels)

    logging.info("Hyper parameter search completed successfully")
    return result


def print_search_results(result):
    """
    Prints all scores and parameters in the hyperparameter search result provided.
    :param result: the result of the hyperparameter search (the output of hyper_param_search/search.fit)
    :return: void
    """
    print("Best: %f using %s" % (result.best_score_, result.best_params_))
    means = result.cv_results_['mean_test_score']
    stds = result.cv_results_['std_test_score']
    params = result.cv_results_['params']
    for mean, stdev, param in zip(means, stds, params):
        print("%f (%f) with: %r" % (mean, stdev, param))
