import logging

from sklearn.neighbors import KNeighborsRegressor

from lib.models import SKModel


class KNNModel(SKModel):

    def __init__(self, route_id, custom_params):
        self.params = {'n_neighbors': 10}
        super().__init__(route_id, custom_params)

    def __create_model__(self):
        logging.info("Creating model...")
        model = KNeighborsRegressor(kwargs=self.params)
        self.model = model
