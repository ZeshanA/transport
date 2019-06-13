import pandas as pd
import sklearn
from sklearn_pandas import gen_features, DataFrameMapper

from lib.date import parse_datetime

TEXT_COLS = ["operator_ref", "progress_rate", "occupancy", "stop_point_ref"]


def convert_journey_to_features(journey):
    df = journey_to_dataframe(journey)
    feature_def = gen_features(
        columns=TEXT_COLS,
        classes=[sklearn.preprocessing.LabelEncoder]
    )
    mapper = DataFrameMapper(feature_def, default=None)
    sample = mapper.fit_transform(df)
    return sample


def journey_to_dataframe(journey):
    timestamp = parse_datetime(journey['Timestamp'])
    eat = journey['ExpectedArrivalTime']
    if eat is None:
        estimate = 240
    else:
        eat = parse_datetime(eat)
        estimate = (eat - timestamp).total_seconds()
    extracted_fields = {
        "direction_ref": [int(journey['DirectionRef'])],
        "operator_ref": [journey['OperatorRef']],
        "longitude": [float(journey['Longitude'])],
        "latitude": [float(journey['Latitude'])],
        "progress_rate": [journey['ProgressRate']],
        "occupancy": [journey['Occupancy']],
        "distance_from_stop": [int(journey['DistanceFromStop'])],
        "stop_point_ref": [journey['StopPointRef']],
        "day": [int(timestamp.strftime('%d'))],
        "month": [int(timestamp.strftime('%m'))],
        "year": [int(timestamp.strftime('%Y'))],
        "hour": [int(timestamp.strftime('%H'))],
        "minute": [int(timestamp.strftime('%M'))],
        "second": [int(timestamp.strftime('%S'))],
        "estimate": [estimate]
    }
    return pd.DataFrame.from_dict(extracted_fields)
