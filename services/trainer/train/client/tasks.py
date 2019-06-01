import asyncio
import logging

import train.client.config as config
from lib.data import get_numpy_datasets
from lib.logs import print_separator
from models import model_types
from param_search.search import hyper_param_search
from train.client.requests import upload_performance_metrics, complete_route, upload_best_parameter_set


def train_route_id(websocket, route_id, model_params=None, train=None, test=None):
    model_class = model_types[config.get_model_type()]
    # Get train/val/test datasets if not provided
    if not train or not test:
        train, test = get_numpy_datasets(route_id, False)
    # Create the requested model
    model = model_class(route_id, **model_params)
    # Train the model
    model.train(train)
    # Calculate and upload final model performance metrics
    metrics = model.calculate_performance_metrics(test)
    asyncio.run(upload_performance_metrics(websocket, metrics))
    # Upload model to object storage
    model.upload_model()
    # Mark model training as complete
    asyncio.run(complete_route(websocket, route_id))
    # Print ASCII divider for clarity in console
    print_separator()


def param_search(websocket, route_id):
    model_class = model_types[config.get_model_type()]
    # Get train/val/test datasets
    train, test = get_numpy_datasets(route_id, validation_set_required=False)
    # Perform hyper parameter search
    result = hyper_param_search(model_class, train)
    # Upload best parameters
    asyncio.run(upload_best_parameter_set(websocket, {'routeID': route_id, **result}))
    # Train final model with the best hyperparameter set
    logging.info("Training final model...")
    # Train and upload the final model
    train_route_id(websocket, route_id, model_params=result['params'], train=train, test=test)


tasks = {
    '-p': param_search,
    '-s': train_route_id
}
