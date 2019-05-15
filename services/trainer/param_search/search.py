from sklearn.model_selection import RandomizedSearchCV
from tensorflow import keras
from tensorflow.keras import layers
from tensorflow.python.keras.wrappers.scikit_learn import KerasRegressor

from lib.args import extract_route_id
from lib.data import COL_COUNT, get_numpy_datasets


def main():
    # Get route_id from CLI arguments
    route_id = extract_route_id()
    # Get train/val/test datasets
    train, test, val = get_numpy_datasets(route_id)
    # Perform hyper parameter search
    hyper_param_search(train, val)


# Perform a randomised grid search
def hyper_param_search(training, validation):
    training_data, training_labels = training

    # Define the type of model we'll be using
    model = KerasRegressor(build_fn=create_model)

    # Define iterable ranges for each of our hyperparameters
    param_dist = {
        'hidden_layer_count': [x for x in range(10, 100)],
        'neuron_count': [x for x in range(256, 1024, 64)],
        'activation_function': ['relu', 'sigmoid', 'tanh'],
        'epochs': [x for x in range(10, 200, 10)]
    }

    # Define the parameters for the search itself
    iteration_count = 1
    random_search = RandomizedSearchCV(estimator=model,
                                       param_distributions=param_dist,
                                       n_iter=iteration_count,
                                       n_jobs=-1,
                                       cv=2,
                                       verbose=3)

    # Perform the search
    result = random_search.fit(training_data, training_labels, validation_data=validation)

    # Display results
    print_search_results(result)


# Build function used by SciKit to create a Keras classifier
def create_model(hidden_layer_count, neuron_count, activation_function):
    # Start constructing a sequential model
    model = keras.Sequential()
    model.add(layers.Dense(COL_COUNT, input_shape=(COL_COUNT,)))

    # Add additional hidden layers as needed
    for i in range(hidden_layer_count - 1):
        model.add(layers.Dense(neuron_count, activation=activation_function))

    # Output layer, a single number
    model.add(layers.Dense(1))

    # Compile model
    model.compile(loss='mean_squared_error',
                  optimizer='adam',
                  metrics=['mean_absolute_error', 'mean_squared_error'])

    return model


def print_search_results(result):
    print("Best: %f using %s" % (result.best_score_, result.best_params_))
    means = result.cv_results_['mean_test_score']
    stds = result.cv_results_['std_test_score']
    params = result.cv_results_['params']
    for mean, stdev, param in zip(means, stds, params):
        print("%f (%f) with: %r" % (mean, stdev, param))


if __name__ == "__main__":
    main()
