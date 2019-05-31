import logging

import numpy as np
from sklearn import svm

from lib.models import SKModel


class SVMModel(SKModel):
    param_dist = {
        'C': [c for c in np.arange(0.5, 50, 0.5)],
        'epsilon': [e for e in np.arange(0.01, 10, 0.02)]
    }

    def __init__(self, route_id, custom_params):
        self.params = {'C': 50, 'epsilon': 5}
        super().__init__(route_id, custom_params)

    @staticmethod
    def create_model(params):
        logging.info("Creating model...")
        model = svm.SVR(C=params['C'], epsilon=params['epsilon'])
        return model
