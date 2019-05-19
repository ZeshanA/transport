import logging
import os
import sys

from sklearn.metrics import mean_absolute_error, mean_squared_error, r2_score
from sklearn.model_selection import RandomizedSearchCV

from tensorflow.python.keras.wrappers.scikit_learn import KerasRegressor

from lib.models import create_model

sys.path.insert(0, os.path.abspath('.'))

from lib.args import extract_cli_args
from lib.data import COL_COUNT, get_numpy_datasets, merge_np_tuples
from lib.files import save_json


def main():
    # Get route_id from CLI arguments
    route_id, base_path = extract_cli_args()
    # Get train/val/test datasets
    train, val, test = get_numpy_datasets(route_id)
    # Perform hyper parameter search
    result = hyper_param_search(train, val)
    # Display results
    print_search_results(result)
    # Save best params
    logging.info("Saving best parameters...")
    save_json(route_id, result.best_params_, base_path, "bestParams.json")
    # Train final model with the best hyperparameter set
    logging.info("Training final model...")
    model = train_final_model(result, merge_np_tuples(train, val))
    # Evaluate final performance and save metric in a file
    logging.info("Saving performance metrics for final model...")
    save_performance_metrics(route_id, model, test, base_path)
    # Save model to disk: disabled for now
    # model.save('models/{}/finalModel.h5'.format(route_id))


def save_performance_metrics(route_id, model, test, base_path):
    """
    Evaluates a model using the test data provided and writes the calculated
    metrics to models/{routeID}/finalPerf.json
    :param route_id: the route id currently being calculated
    :param model: the Keras model to evaluate (any model with support for .predict() should work)
    :param test: a pair of Numpy arrays in the format (testing_data, testing_labels)
    :return:
    """
    data, labels = test
    preds = model.predict(data)
    metrics = {
        'mean_absolute_error': mean_absolute_error(labels, preds),
        'mean_squared_error': mean_squared_error(labels, preds),
        'r2_score': r2_score(labels, preds)
    }
    save_json(route_id, metrics, base_path, "finalPerf.json")


def hyper_param_search(training, validation):
    """
    Performs a randomised grid search using the given training and validation data sets.
    :param training: a pair of Numpy arrays in the format (training_data, training_labels)
    :param validation: a pair of Numpy arrays in the format (validation_data, validation_labels)
    :return: SciKit.cv_results_ object containing the results of the search
    """
    logging.info("Starting hyper parameter search...")
    training_data, training_labels = training

    # Define the type of model we'll be using
    model = KerasRegressor(build_fn=create_model)

    # Define iterable ranges for each of our hyperparameters
    param_dist = {
        'hidden_layer_count': [x for x in range(10, 100)],
        'neuron_count': [x for x in range(256, 1024, 64)],
        'activation_function': ['relu', 'tanh'],
        'epochs': [x for x in range(10, 40, 4)],
    }

    # Define the parameters for the search itself
    iteration_count = 4
    random_search = RandomizedSearchCV(estimator=model,
                                       param_distributions=param_dist,
                                       n_iter=iteration_count,
                                       n_jobs=-1,
                                       cv=2,
                                       verbose=3)

    # Perform the search and return the results
    result = random_search.fit(training_data, training_labels, validation_data=validation)

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


def train_final_model(result, training):
    """
    Trains a new model using the best parameters from a hyperparameter search.
    :param result: the result of the hyperparameter search (the output of hyper_param_search/search.fit)
    :param training: a pair of Numpy arrays in the format (training_data, training_labels)
    :return: a trained Keras model backed by Tensorflow
    """
    best_params = result.best_params_
    model = create_model(best_params['hidden_layer_count'], best_params['neuron_count'],
                         best_params['activation_function'])
    training_data, training_labels = training
    model.fit(x=training_data, y=training_labels, epochs=best_params['epochs'])
    return model


if __name__ == "__main__":
    logging.basicConfig()
    logging.getLogger().setLevel(logging.INFO)
    main()
