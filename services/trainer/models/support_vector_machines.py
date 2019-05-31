import logging

from sklearn import svm

from lib.models import SKModel


class SVMModel(SKModel):

    def __init__(self, route_id, custom_params):
        self.params = {'C': 50, 'epsilon': 5}
        super().__init__(route_id, custom_params)

    def __create_model__(self):
        logging.info("Creating model...")
        model = svm.SVR(C=self.params['C'], epsilon=self.params['epsilon'])
        self.model = model
