import logging

from sklearn.tree import DecisionTreeRegressor

from lib.models import SKModel


class DecisionTreeModel(SKModel):
    param_dist = {}

    def __init__(self, route_id, custom_params):
        self.params = {}
        super().__init__(route_id, custom_params)

    @staticmethod
    def create_model(self):
        logging.info("Creating model...")
        return DecisionTreeRegressor()
