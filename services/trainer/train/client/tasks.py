import asyncio
import logging

from lib.data import get_numpy_datasets, merge_np_tuples
from lib.logs import print_separator
from models import model_types
from param_search.search import hyper_param_search, print_search_results
import train.client.config as config
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


async def param_search(websocket, route_id):
    loop = asyncio.get_event_loop()
    model_class = model_types[config.get_model_type()]
    # Get train/val/test datasets
    train, val, test = await loop.run_in_executor(None, get_numpy_datasets, route_id)
    # Perform hyper parameter search
    result = await loop.run_in_executor(None, hyper_param_search, model_class, train, val)
    # Display results
    print_search_results(result)
    # Upload best params
    # await upload_best_parameter_set(websocket, result.best_params_)
    # Train final model with the best hyperparameter set
    logging.info("Training final model...")
    final_training_set = merge_np_tuples(train, val)
    await train_route_id(websocket, route_id, train=final_training_set, test=test)


tasks = {
    '-p': param_search,
    '-s': train_route_id
}
