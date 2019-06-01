import logging

import matplotlib.pyplot as plt
import pandas as pd
from skopt.space import Integer, Categorical
from tensorflow import keras
# noinspection PyUnresolvedReferences
from tensorflow.keras import layers

from lib.data import COL_COUNT
from lib.models import Model


class NNModel(Model):
    param_dist = [
        Integer(10, 100, name='hidden_layer_count'),
        Integer(256, 1024, name='neuron_count'),
        Categorical(['relu', 'tanh'], name='activation_function'),
        Integer(10, 40, name='epochs')
    ]
    default_params = {
        'hidden_layer_count': 47,
        'neuron_count': 552,
        'activation_function': 'relu',
        'epochs': 18
    }

    @staticmethod
    def create_model(
            hidden_layer_count=default_params['hidden_layer_count'],
            neuron_count=default_params['neuron_count'],
            activation_function=default_params['activation_function'],
            epochs=default_params['epochs']
    ):
        # Start constructing a sequential model
        model = keras.Sequential()
        model.add(layers.Dense(COL_COUNT, input_shape=(COL_COUNT,)))

        # Add additional hidden layers as needed
        for i in range(hidden_layer_count - 1):
            model.add(layers.Dense(neuron_count, activation=activation_function))

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
