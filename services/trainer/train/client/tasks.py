import asyncio
import logging

import train.client.config as config
from lib.data import get_numpy_datasets
from lib.logs import print_separator
from models import model_types
from param_search.search import hyper_param_search
from train.client.requests import upload_performance_metrics, complete_route, upload_best_parameter_set


async def train_route_id(websocket, route_id, train=None, test=None):
    loop = asyncio.get_event_loop()
    model_class = model_types[config.get_model_type()]
    # Get train/val/test datasets if not provided
    if not train or not test:
        train, test = await loop.run_in_executor(None, get_numpy_datasets, route_id, False)
    # Create the requested model
    model = model_class(route_id)
    # Train the model
    await loop.run_in_executor(None, model.train, train)
    # Calculate and upload final model performance metrics
    metrics = await loop.run_in_executor(None, model.calculate_performance_metrics, test)
    await upload_performance_metrics(websocket, metrics)
    # Upload model to object storage
    await loop.run_in_executor(None, model.upload_model)
    # Mark model training as complete
    await complete_route(websocket, route_id)
    # Print ASCII divider for clarity in console
    print_separator()


def param_search(websocket, route_id):
    model_class = model_types[config.get_model_type()]
    # Get train/val/test datasets
    train, test = get_numpy_datasets(route_id, validation_set_required=False)
    # Perform hyper parameter search
    result = hyper_param_search(model_class, train)
    # Upload best params
    upload_best_parameter_set(websocket, result)
    # Train final model with the best hyperparameter set
    logging.info("Training final model...")
    train_route_id(websocket, route_id, train=train, test=test)


tasks = {
    '-p': param_search,
    '-s': train_route_id
}
