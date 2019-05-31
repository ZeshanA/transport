from models.k_nearest_neighbours import KNNModel
from models.neural_network import NNModel
from models.random_forest import RandomForestModel
from models.support_vector_machines import SVMModel

model_types = {
    'neural_network': NNModel,
    'random_forest': RandomForestModel,
    'support_vector_machine': SVMModel,
    'k_nearest_neighbours': KNNModel
}
