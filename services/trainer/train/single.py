import sys

from lib.args import get_model_type
from lib.data import get_numpy_datasets
from lib.logs import init_logging
from models.neural_network import NNModel
from models.random_forest import RandomForestModel

MODEL_TYPES = {
    'neural_network': NNModel,
    'random_forest': RandomForestModel
}


def train_model():
    model_class, route_id = MODEL_TYPES[get_model_type()], sys.argv[2]
    # Get train/val/test datasets
    train, test = get_numpy_datasets(route_id, False)
    # Create the requested model
    model = model_class(route_id)
    # Train the model
    model.train(train)
    # Calculate and print final model performance metrics
    metrics = model.calculate_performance_metrics(test)
    print(metrics)


if __name__ == '__main__':
    init_logging()
    train_model()
