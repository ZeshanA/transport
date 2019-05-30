from lib.args import get_model_type
from lib.data import get_numpy_datasets
from lib.logs import print_separator
from models.neural_network import NNModel
from models.random_forest import RandomForestModel
from train.client.requests import upload_performance_metrics, complete_route

MODEL_TYPES = {
    'neural_network': NNModel,
    'random_forest': RandomForestModel
}


async def train_route_id(websocket, route_id):
    model_class = MODEL_TYPES[get_model_type()]
    # Get train/val/test datasets
    train, test = get_numpy_datasets(route_id, False)
    # Create the requested model
    model = model_class(route_id)
    # Train the model
    model.train(train)
    # Calculate and upload final model performance metrics
    metrics = model.calculate_performance_metrics(test)
    await upload_performance_metrics(websocket, metrics)
    # Upload model to object storage
    model.upload_model()
    # Mark model training as complete
    await complete_route(websocket, route_id)
    # Print ASCII divider for clarity in console
    print_separator()
