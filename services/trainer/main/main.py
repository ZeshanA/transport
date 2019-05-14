import logging
import os
import sys

import matplotlib.pyplot as plt
import pandas as pd
import tensorflow as tf
from sklearn.model_selection import train_test_split
from tensorflow import keras, feature_column
from tensorflow.keras import layers

from db.db import connect

LABEL_COL = "time_to_stop"
NUMERIC_COLS = ["direction_ref", "longitude", "latitude", "distance_from_stop",
                "day", "month", "year", "hour", "minute", "second", "estimate"]
TEXT_COLS = ["operator_ref", "progress_rate", "occupancy", "stop_point_ref"]


def main():
    logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(levelname)s: %(message)s')
    # Get route_id from CLI arguments
    route_id = extract_route_id()
    # Get dataframe from DB
    train, val, test = get_datasets(route_id, LABEL_COL)
    # Create the model's structure
    model = get_model()
    # Train the model on training data and save in checkpoints as you go
    train_model(model, train, val, [get_checkpoint_callback(route_id)])


# Extracts the route_id from the CLI arguments, exiting if none was provided
def extract_route_id():
    if len(sys.argv) < 2:
        logging.critical("No route ID was provided, correct usage: ./trainer <routeID>")
        sys.exit(1)
    return sys.argv[1]


# Creates a folder for this route_id and returns a keras callback that can
# be passed to model.fit to save the model to the folder after every epoch
def get_checkpoint_callback(route_id):
    os.makedirs("models/{}".format(route_id))
    checkpoint_path = "models/{}/cp.ckpt".format(route_id)
    cp_callback = keras.callbacks.ModelCheckpoint(checkpoint_path)
    return cp_callback


def get_model():
    feature_cols = get_feature_columns(NUMERIC_COLS, TEXT_COLS)
    feature_layer = layers.DenseFeatures(feature_cols)
    model = keras.Sequential([
        feature_layer,
        layers.Dense(512, activation='relu'),
        layers.Dense(512, activation='relu'),
        layers.Dense(512, activation='relu'),
        layers.Dense(512, activation='relu'),
        layers.Dense(1)
    ])
    metrics = ['mean_absolute_error', 'mean_squared_error']
    model.compile(loss='mean_squared_error', optimizer='adam', metrics=metrics)
    return model


def train_model(model, train, validation, callbacks):
    history = model.fit(train, validation_data=validation, epochs=50, callbacks=callbacks)
    plot_history(history)


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


# Fetches data for route_id, splits it into (train, test, val) and returns
# each set as a tf.data.Dataset
def get_datasets(route_id: str, label_col: str, batch_size: int = 32):
    logging.info('Fetching data for route_id: {}'.format(route_id))
    dataframe = get_dataframe(route_id)

    pd.set_option('display.max_columns', 500)
    print(dataframe[dataframe.isnull().T.any().T])

    # Split into train/val/test sets
    train, test = train_test_split(dataframe, test_size=0.2)
    train, val = train_test_split(train, test_size=0.2)
    # Convert each set into tf.data format
    train_ds = df_to_dataset(train, label_col, batch_size=batch_size)
    val_ds = df_to_dataset(val, label_col, batch_size=batch_size)
    test_ds = df_to_dataset(test, label_col, batch_size=batch_size)
    logging.info("Succesfully fetched, split and converted data for route_id: {}".format(route_id))
    return train_ds, val_ds, test_ds


# Fetches the data from the DB for the current route_id
# and returns it in a Pandas dataframe
def get_dataframe(route_id: str) -> pd.DataFrame:
    conn = connect()
    dataframe = pd.read_sql(
        """
        SELECT 
            direction_ref, operator_ref, longitude, latitude, progress_rate,
            COALESCE(occupancy, '') AS occupancy, distance_from_stop, stop_point_ref,
            EXTRACT(day from timestamp) AS day, EXTRACT(month from timestamp) AS month, EXTRACT(year from timestamp) AS year,
            EXTRACT(hour from timestamp) AS hour, EXTRACT(minute from timestamp) as minute, EXTRACT(second from timestamp) as second,
            COALESCE(EXTRACT(epoch FROM expected_arrival_time - timestamp)::integer, 0) AS estimate,
            time_to_stop
        FROM labelled_journey
        WHERE line_ref='{}';
        """.format(route_id),
        conn
    )
    return dataframe


# Converts a pandas dataframe to tf.data format
def df_to_dataset(dataframe: pd.DataFrame, label_col: str, batch_size: int):
    dataframe = dataframe.copy()
    labels = dataframe.pop(label_col)
    ds = tf.data.Dataset.from_tensor_slices((dict(dataframe), labels))
    ds = ds.shuffle(buffer_size=len(dataframe))
    ds = ds.batch(batch_size)
    return ds


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


if __name__ == "__main__":
    main()
