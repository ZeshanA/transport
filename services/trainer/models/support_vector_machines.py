import logging

from sklearn import svm
from skopt.space import Real

from lib.models import SKModel


class SVMModel(SKModel):
    param_dist = [
        Real(-3, 2, name='C'),
        Real(-3, 2, name='epsilon')
    ]
    default_params = {'C': 1, 'epsilon': 0.1}

    @staticmethod
    def create_model(c=default_params['C'], epsilon=default_params['epsilon']):
        logging.info("Creating model...")
        return svm.SVR(C=c, epsilon=epsilon)
