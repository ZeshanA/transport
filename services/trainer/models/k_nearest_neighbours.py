import logging

from sklearn.neighbors import KNeighborsRegressor
from skopt.space import Integer, Categorical

from lib.models import SKModel


class KNNModel(SKModel):
    param_dist = [
        Integer(5, 1000, name='n_neighbors'),
        Categorical(name='weights', categories=['uniform', 'distance']),
        Categorical(name='metric', categories=['euclidean', 'manhattan', 'chebyshev', 'minkowski']),
    ]
    default_params = {
        'n_neighbors': 10,
        'weights': 'uniform',
        'metric': 'minkowski'
    }

    @staticmethod
    def create_model(n_neighbors=default_params['n_neighbors'], weights=default_params['weights'],
                     metric=default_params['metric']):
        logging.info("Creating model...")
        return KNeighborsRegressor(n_neighbors=n_neighbors, weights=weights, metric=metric)
