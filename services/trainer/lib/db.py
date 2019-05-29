import logging
import os
import psycopg2

HOST = "mtadata.postgres.database.azure.com"
PORT = 5432
NAME = "postgres"
TIME_FORMAT = "2006-01-02 15:04:05"
DATE_FORMAT = "2006-01-02"


def connect():
    username, password = get_db_details()
    conn = psycopg2.connect(host=HOST, database=NAME, user=username, password=password)
    return conn


def get_db_details():
    username = os.environ.get('TRANSPORT_DB_USERNAME')
    password = os.environ.get('TRANSPORT_DB_PASSWORD')
    if not username or not password:
        logging.critical('TRANSPORT_DB_{USERNAME/PASSWORD} NOT SET')
        raise KeyError
    return username, password
