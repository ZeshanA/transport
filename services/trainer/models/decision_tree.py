import logging

import joblib
from sklearn.tree import DecisionTreeRegressor

from lib.models import Model


class DecisionTreeModel(Model):

    def __init__(self, route_id):
        super().__init__(route_id, {'n_neigbours': 10})

    def __create_model__(self):
        logging.info("Creating model...")
        model = DecisionTreeRegressor()
        self.model = model

    def train(self, training):
        data, labels = training
        logging.info("Training model...")
        self.model.fit(data, labels)
        logging.info("Model successfully trained...")

    def __save_model__(self, filepath):
        joblib.dump(self.model, filepath)
