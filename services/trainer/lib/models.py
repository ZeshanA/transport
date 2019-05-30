import logging
import os
from abc import ABC, abstractmethod

import boto3
from sklearn.metrics import mean_absolute_error, mean_squared_error, r2_score

from lib.files import save_json
from lib.storage import get_storage_details


class Model(ABC):
    def __init__(self, route_id, params):
        self.route_id, self.params = route_id, params
        self.model, self.history = None, None
        self.model_name = type(self).__name__
        self.__create_model__()

    @abstractmethod
    def __create_model__(self):
        pass

    @abstractmethod
    def train(self, training):
        pass

    @abstractmethod
    def __save_model__(self, filepath):
        pass

    def save_model_to_disk(self):
        """
        Saves a trained Tensorflow model to disk. Can be loaded again using Keras'
        load_model function.
        :param route_id: string: routeID for the current model
        :param model: pointer to a trained Tensorflow model
        :return: the (relative) filepath that the model was saved at
        """

        directory = '/data/za816/trained/{}/{}/'.format(self.model_name, self.route_id)
        filepath = '{}/finalModel.h5'.format(directory, self.route_id)
        os.makedirs(directory, exist_ok=True)
        self.__save_model__(filepath)
        return filepath

    def calculate_performance_metrics(self, test):
        logging.info("Calculating model performance metrics for routeID %s...", self.route_id)
        data, labels = test
        preds = self.model.predict(data)
        return {
            'route_id': self.route_id,
            'mean_absolute_error': mean_absolute_error(labels, preds),
            'mean_squared_error': mean_squared_error(labels, preds),
            'r2_score': r2_score(labels, preds)
        }

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
        logging.info("Successfully saved model performance metrics for routeID %s...", self.route_id)

    def upload_model(self):
        """
        Saves the model to cloud object storage.
        """
        logging.info("Uploading final model for routeID %s to storage...", self.route_id)
        filepath = self.save_model_to_disk()
        session = boto3.session.Session()
        key_id, secret = get_storage_details()
        client = session.client('s3',
                                region_name='fra1',
                                endpoint_url='https://fra1.digitaloceanspaces.com',
                                aws_access_key_id=key_id,
                                aws_secret_access_key=secret)
        client.upload_file(filepath, 'mtadata3', '{}-{}-finalModel.h5'.format(self.model_name, self.route_id),
                           ExtraArgs={'ACL': 'public-read'})
        logging.info("Successfully uploaded final model for routeID %s to storage...", self.route_id)
