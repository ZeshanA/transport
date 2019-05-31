import logging

from sklearn.tree import DecisionTreeRegressor

from lib.models import SKModel


class DecisionTreeModel(SKModel):

    def __init__(self, route_id, custom_params):
        self.params = {}
        super().__init__(route_id, custom_params)

    def __create_model__(self):
        logging.info("Creating model...")
        model = DecisionTreeRegressor()
        self.model = model
