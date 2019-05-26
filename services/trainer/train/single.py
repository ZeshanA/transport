import logging
import os
import sys

sys.path.insert(0, os.path.abspath('.'))

from lib.data import get_numpy_datasets, merge_np_tuples
from lib.logs import init_logging
from lib.models import create_model, calculate_performance_metrics

OPTIMAL_PARAMS = {
    'hidden_layer_count': 47,
    'neuron_count': 552,
    'activation_function': 'relu',
    'epochs': 18
}


def main():
    # Fetch next routeID from server
    route_id = sys.argv[1]
    # Get train/val/test datasets
    train, val, test = get_numpy_datasets(route_id)
    # Train model using pre-determined optimal parameters
    model = get_trained_model(OPTIMAL_PARAMS, merge_np_tuples(train, val))
    # Calculate and upload final model performance metrics
    metrics = calculate_performance_metrics(route_id, model, test)
    print(metrics)


def get_trained_model(params, training):
    """
    Creates a model based on the given parameters, trains it on the provided training data
    and returns a pointer to the fully trained model.
    :param params: a dict containing the parameters that the model should be trained using
    :param training: a tuple of Numpy arrays (training_data, training_labels)
    :return: a pointer to a trained Tensorflow model
    """
    logging.info("Starting model training...")
    training_data, training_labels = training
    model = create_model(params['hidden_layer_count'], params['neuron_count'], params['activation_function'])
    model.fit(x=training_data, y=training_labels, epochs=params['epochs'])
    logging.info("Successfully completed model training...")
    return model


def print_separator():
    """
    Prints an ASCII divider to separate different executions in the console
    """
    print("\n\n===================================================================================================\n\n")


if __name__ == "__main__":
    init_logging()
    main()
