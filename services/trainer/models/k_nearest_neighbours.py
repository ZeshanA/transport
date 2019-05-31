import logging

from sklearn.neighbors import KNeighborsRegressor

from lib.models import SKModel


class KNNModel(SKModel):
    param_dist = {'n_neighbors': [x for x in range(2, 50, 2)]}

    def __init__(self, route_id, custom_params):
        self.params = {'n_neighbors': 10}
        super().__init__(route_id, custom_params)

    @staticmethod
    def create_model(n_neighbors):
        logging.info("Creating model...")
        return KNeighborsRegressor(n_neighbors=n_neighbors)

