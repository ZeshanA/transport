import os

import matplotlib.pyplot as plt
import pandas as pd
from tensorflow import keras, feature_column
from tensorflow.keras import layers

from lib.data import COL_COUNT


def create_model(hidden_layer_count, neuron_count, activation_function):
    """
    Build function used by SciKit to create a Keras classifier.
    :param hidden_layer_count: number of intermediary layers in the network (excluding feature layer)
    :param neuron_count: number of neurons in each layer
    :param activation_function: the activation function applied by each neuron
    :return: an untrained Keras model backed by Tensorflow
    """
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


# Returns a list of tensorflow feature columns, converting text_cols into
# categorical/embedded columns as needed
def get_feature_columns(numeric_cols, text_cols):
    feature_cols = []
    for col in numeric_cols:
        feature_cols.append(feature_column.numeric_column(col))
    for col in text_cols:
        cat_col = feature_column.categorical_column_with_hash_bucket(col, 10000)
        emb_col = feature_column.embedding_column(cat_col, 8)
        feature_cols.append(emb_col)
    # TODO: maybe remove this line
    print(list(map(lambda x: x.key if hasattr(x, 'key') else x.categorical_column.key, feature_cols)))
    return feature_cols


# Creates a folder for this route_id and returns a keras callback that can
# be passed to model.fit to save the model to the folder after every epoch
def get_checkpoint_callback(route_id):
    os.makedirs("models/{}".format(route_id))
    checkpoint_path = "models/{}/cp.ckpt".format(route_id)
    cp_callback = keras.callbacks.ModelCheckpoint(checkpoint_path)
    return cp_callback


# Plot train error against validation error using a history.
# `history` is the result of calling model.fit()
def plot_history(history):
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
