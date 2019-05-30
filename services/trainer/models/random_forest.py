import logging

import joblib
from sklearn.ensemble import RandomForestRegressor

from lib.models import Model


class RandomForestModel(Model):

    def __init__(self, route_id):
        super().__init__(route_id, {'n_estimators': 200})

    def __create_model__(self):
        logging.info("Creating model...")
        model = RandomForestRegressor(
            n_estimators=self.params['n_estimators'],
            n_jobs=-1
        )
        self.model = model

    def train(self, training):
        data, labels = training
        logging.info("Training model...")
        self.model.fit(data, labels)
        logging.info("Model successfully trained...")

    def __save_model__(self, filepath):
        joblib.dump(self.model, filepath)
