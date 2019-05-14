import logging

import pandas as pd
from sklearn.model_selection import train_test_split
import tensorflow as tf

from db.db import connect


LABEL_COL = "time_to_stop"
NUMERIC_COLS = ["direction_ref", "longitude", "latitude", "distance_from_stop",
                "day", "month", "year", "hour", "minute", "second", "estimate"]
TEXT_COLS = ["operator_ref", "progress_rate", "occupancy", "stop_point_ref"]


# Fetches data for route_id, splits it into (train, test, val) and returns
# each set as a tf.data.Dataset
def get_datasets(route_id: str, batch_size: int = 32):
    logging.info('Fetching data for route_id: {}'.format(route_id))
    dataframe = get_dataframe(route_id)

    pd.set_option('display.max_columns', 500)
    print(dataframe[dataframe.isnull().T.any().T])

    # Split into train/val/test sets
    train, test = train_test_split(dataframe, test_size=0.2)
    train, val = train_test_split(train, test_size=0.2)
    # Convert each set into tf.data format
    train_ds = df_to_dataset(train, LABEL_COL, batch_size=batch_size)
    val_ds = df_to_dataset(val, LABEL_COL, batch_size=batch_size)
    test_ds = df_to_dataset(test, LABEL_COL, batch_size=batch_size)
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
