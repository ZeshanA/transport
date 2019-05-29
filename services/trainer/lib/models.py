import logging
from abc import ABC, abstractmethod

from lib.files import save_json


class Model(ABC):
    def __init__(self, route_id, params):
        self.route_id, self.params = route_id, params
        self.model, self.history = None, None

    @abstractmethod
    def create_model(self):
        pass

    @abstractmethod
    def train(self, training):
        pass

    @abstractmethod
    def calculate_performance_metrics(self, test):
        pass

    def save_performance_metrics(self, test, base_path):
        """
        Evaluates a model using the test data provided and writes the calculated
        metrics to models/{routeID}/finalPerf.json
        :param test: a pair of Numpy arrays in the format (testing_data, testing_labels)
        :param base_path: the base_path under which to save the performance data
        :return:
        """
        metrics = self.calculate_performance_metrics(test)
        save_json(self.route_id, metrics, base_path, "finalPerf.json")
        logging.info("Successfully saved model performance metrics for routeID %s...", route_id)

