import logging

from tensorflow import keras
from tensorflow.keras import layers

from lib.args import extract_route_id
from lib.data import get_datasets, COL_COUNT, NUMERIC_COLS, TEXT_COLS
from lib.models import plot_history, get_checkpoint_callback, get_feature_columns


def main():
    logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(levelname)s: %(message)s')
    # Get route_id from CLI arguments
    route_id = extract_route_id()
    # Get training, validation, testing datasets
    train, val, test = get_datasets(route_id)
    # Create the model's structure
    model = create_model(14, 320, 'tanh')
    # Train the model on training data and save in checkpoints as you go
    train_model(model, train, val, [get_checkpoint_callback(route_id)])


# Build function used by SciKit to create a Keras classifier
def create_model(hidden_layer_count, neuron_count, activation_function):
    # Start constructing a sequential model
    model = keras.Sequential()
    feature_cols = get_feature_columns(NUMERIC_COLS, TEXT_COLS)
    feature_layer = layers.DenseFeatures(feature_cols)
    model.add(feature_layer)

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


def train_model(model, train, validation, callbacks):
    history = model.fit(train, validation_data=validation, epochs=125, callbacks=callbacks)
    plot_history(history)


if __name__ == "__main__":
    main()
