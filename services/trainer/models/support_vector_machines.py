import logging

from sklearn import svm
from skopt.space import Real

from lib.models import SKModel


class SVMModel(SKModel):
    param_dist = [
        Real(0.01, 50, name='C'),
        Real(0.001, 10, name='epsilon')
    ]
    default_params = {'C': 50, 'epsilon': 5}

    @staticmethod
    def create_model(c=default_params['C'], epsilon=default_params['epsilon']):
        logging.info("Creating model...")
        return svm.SVR(C=c, epsilon=epsilon)
