import logging

from tensorflow import keras
from tensorflow.keras import layers

from lib.args import extract_route_id
from lib.data import get_datasets, NUMERIC_COLS, TEXT_COLS
from lib.models import plot_history, get_feature_columns, get_checkpoint_callback


def main():
    logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(levelname)s: %(message)s')
    # Get route_id from CLI arguments
    route_id = extract_route_id()
    # Get dataframe from DB
    train, val, test = get_datasets(route_id)
    # Create the model's structure
    model = get_model()
    # Train the model on training data and save in checkpoints as you go
    train_model(model, train, val, [get_checkpoint_callback(route_id)])


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


if __name__ == "__main__":
    main()
