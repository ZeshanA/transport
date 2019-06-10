import logging

import gc
import numpy as np
import pandas as pd
import sklearn
from sklearn.model_selection import train_test_split
from sklearn_pandas import gen_features, DataFrameMapper

from lib.db import connect

LABEL_COL = "time_to_stop"
NUMERIC_COLS = ["direction_ref", "longitude", "latitude", "distance_from_stop",
                "day", "month", "year", "hour", "minute", "second", "estimate"]
TEXT_COLS = ["operator_ref", "progress_rate", "occupancy", "stop_point_ref"]
COL_COUNT = len(NUMERIC_COLS) + len(TEXT_COLS)


def get_numpy_datasets(route_id: str, validation_set_required: bool = True):
    df = get_dataframe(route_id)
    train, test = train_test_split(df, test_size=0.2)
    gc.collect()  # Attempt to reduce memory consumption from dataframe
    train, val = train_test_split(train, test_size=0.2)
    feature_def = gen_features(
        columns=TEXT_COLS,
        classes=[sklearn.preprocessing.LabelEncoder]
    )
    mapper = DataFrameMapper(feature_def, default=None)
    train_labels, val_labels, test_labels = train.pop(LABEL_COL), val.pop(LABEL_COL), test.pop(LABEL_COL)
    train_data, val_data, test_data = mapper.fit_transform(train), mapper.fit_transform(val), mapper.fit_transform(test)
    train, val, test = (train_data, train_labels), (val_data, val_labels), (test_data, test_labels)
    logging.info("Successfully split and converted data for route_id {}".format(route_id))
    if not validation_set_required:
        return merge_np_tuples(train, val), test
    return train, val, test


# Fetches the data from the DB for the current route_id
# and returns it in a Pandas dataframe
def get_dataframe(route_id: str) -> pd.DataFrame:
    logging.info("Fetching data for route_id {}".format(route_id))
    conn = connect()
    dataframe = pd.read_sql(
        """
        SELECT 
            direction_ref, operator_ref, longitude, latitude, progress_rate,
            COALESCE(occupancy, '') AS occupancy, distance_from_stop, stop_point_ref,
            EXTRACT(day from timestamp) AS day,
            EXTRACT(month from timestamp) AS month, 
            EXTRACT(year from timestamp) AS year,
            EXTRACT(hour from timestamp) AS hour,
            EXTRACT(minute from timestamp) as minute,
            EXTRACT(second from timestamp) as second,
            COALESCE(EXTRACT(epoch FROM expected_arrival_time - timestamp)::integer, 0) AS estimate,
            time_to_stop
        FROM labelled_journey
        WHERE line_ref='{}' LIMIT 100;
        """.format(route_id),
        conn
    )
    logging.info("Successfully fetched data for route_id {}".format(route_id))
    return dataframe


def merge_np_tuples(a, b):
    """
    Takes two tuples containing Numpy arrays and returns a single tuple of the same length,
    where each index contains the concatenation of the arrays at that index in the two input tuples.
    example: merge_tuples(([a,b],[c,d]), ([e,f],[g,h])) = ([a,b,e,f], [c,d,g,h])
    :param a: tuple of numpy arrays
    :param b: tuple of numpy arrays
    :return: a tuple
    """
    result = []
    for item_a, item_b in zip(a, b):
        result.append(np.concatenate((item_a, item_b)))
    return tuple(result)
