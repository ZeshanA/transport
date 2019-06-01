import logging

from sklearn.neighbors import KNeighborsRegressor
from skopt.space import Integer

from lib.models import SKModel


class KNNModel(SKModel):
    param_dist = [Integer(5, 1000, name='n_neighbors')]
    default_params = {'n_neighbors': 10}

    @staticmethod
    def create_model(n_neighbors=default_params['n_neighbors']):
        logging.info("Creating model...")
        return KNeighborsRegressor(n_neighbors=n_neighbors)
