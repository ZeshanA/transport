from datetime import datetime

from pytz import timezone

FORMAT_STRING = '%Y-%m-%d %H:%M:%S'


def parse_datetime(string):
    est = timezone('US/Eastern')
    naive = datetime.strptime(string, FORMAT_STRING)
    return est.localize(naive)


def format_datetime(dt):
    return datetime.strftime(dt, FORMAT_STRING)
