import logging

from sklearn.ensemble import RandomForestRegressor
from skopt.space import Integer

from lib.models import SKModel


class RandomForestModel(SKModel):
    param_dist = [Integer(5, 200, name='n_estimators')]
    default_params = {'n_estimators': 200}

    @staticmethod
    def create_model(n_estimators=default_params['n_estimators']):
        logging.info("Creating model...")
        return RandomForestRegressor(
            n_estimators=n_estimators,
            n_jobs=1
        )
