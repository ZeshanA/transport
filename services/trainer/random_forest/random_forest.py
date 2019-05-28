# Import the model we are using
import logging

from sklearn.ensemble import RandomForestRegressor
import numpy as np

from lib.data import get_numpy_datasets
from lib.logs import init_logging
from lib.routes import all_routes


def main():
    init_logging()
    errors = []
    for route_id in all_routes[:10]:
        errors.append(train_for_route(route_id))
    print("Routes:", all_routes[:10])
    print("Errors:", errors)


def train_for_route(route_id):
    logging.info("Fetching datasets...")
    # Get train/val/test datasets
    train, val, test = get_numpy_datasets(route_id)
    logging.info("Creating model...")
    # Instantiate model with 1000 decision trees
    rf = RandomForestRegressor(n_estimators=200)
    logging.info("Training model...")
    rf.fit(train[0], train[1])
    logging.info("Successfully trained...")
    logging.info("Making predictions...")
    # Use the forest's predict method on the test data
    predictions = rf.predict(test[0])
    logging.info("Calculating errors...")
    # Calculate the absolute errors
    errors = abs(predictions - test[1])
    # Print out the mean absolute error (mae)
    print('Mean Absolute Error: ', np.mean(errors))
    return np.mean(errors)


if __name__ == '__main__':
    main()
