import logging

from sklearn.tree import DecisionTreeRegressor

from lib.models import SKModel


class DecisionTreeModel(SKModel):
    param_dist, default_params = [], {}

    @staticmethod
    def create_model():
        logging.info("Creating model...")
        return DecisionTreeRegressor()
