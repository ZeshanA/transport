import logging

import matplotlib.pyplot as plt
import pandas as pd
from tensorflow import keras
# noinspection PyUnresolvedReferences
from tensorflow.keras import layers

from lib.data import COL_COUNT
from lib.models import Model


class NNModel(Model):
    param_dist = {
        'hidden_layer_count': [x for x in range(10, 100)],
        'neuron_count': [x for x in range(256, 1024, 64)],
        'activation_function': ['relu', 'tanh'],
        'epochs': [x for x in range(10, 40, 4)],
    }

    def __init__(self, route_id, custom_params):
        self.params = {
            'hidden_layer_count': 1,
            'neuron_count': 1,
            'activation_function': 'relu',
            'epochs': 1
        }
        super().__init__(route_id, custom_params)

    @staticmethod
    def create_model(params):
        # Start constructing a sequential model
        model = keras.Sequential()
        model.add(layers.Dense(COL_COUNT, input_shape=(COL_COUNT,)))

        # Add additional hidden layers as needed
        for i in range(params['hidden_layer_count'] - 1):
            model.add(layers.Dense(params['neuron_count'], activation=params['activation_function']))

        # Output layer, a single number
        model.add(layers.Dense(1))

        # Compile model
        model.compile(loss='mean_squared_error',
                      optimizer='adam',
                      metrics=['mean_absolute_error', 'mean_squared_error'])

        return model

    def train(self, training):
        logging.info("Starting model training...")
        training_data, training_labels = training
        self.history = self.model.fit(x=training_data, y=training_labels, epochs=self.params['epochs'])
        logging.info("Successfully completed model training...")

    def __save_model__(self, filepath):
        self.model.save(filepath)

    def plot_history(self):
        """
        Plot train error against validation error using a history.
        `history` is the result of calling model.fit()
        :return:
        """
        history = self.history
        hist = pd.DataFrame(history.history)
        hist['epoch'] = history.epoch

        plt.figure()
        plt.xlabel('Epoch')
        plt.ylabel('Mean Abs Error [MPG]')
        plt.plot(hist['epoch'], hist['mean_absolute_error'], label='Train Error')
        plt.plot(hist['epoch'], hist['val_mean_absolute_error'], label='Val Error')
        plt.legend()

        plt.figure()
        plt.xlabel('Epoch')
        plt.ylabel('Mean Square Error [$MPG^2$]')
        plt.plot(hist['epoch'], hist['mean_squared_error'],
                 label='Train Error')
        plt.plot(hist['epoch'], hist['val_mean_squared_error'],
                 label='Val Error')
        plt.legend()
        plt.show()
