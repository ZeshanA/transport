import logging

from sklearn import svm
from skopt.space import Real, Categorical

from lib.models import SKModel


class SVMModel(SKModel):
    param_dist = [
        Categorical(name='kernel', categories=['rbf', 'poly', 'sigmoid']),
        Real(0.1, 2, name='C'),
        Real(0.1, 2, name='epsilon')
    ]
    default_params = {'kernel': 'rbf', 'C': 1, 'epsilon': 0.1}

    # noinspection PyPep8Naming
    @staticmethod
    def create_model(kernel=default_params['kernel'], C=default_params['C'], epsilon=default_params['epsilon']):
        logging.info("Creating model...")
        return svm.SVR(kernel=kernel, C=C, epsilon=epsilon, gamma='auto')
