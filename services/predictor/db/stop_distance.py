from db.db import connect

AVG_DISTANCE = 485


# TODO: Clean up this function
def get_stop_distance(routeID, directionID, fromStop, toStop):
    conn = connect()
    cur = conn.cursor()
    cur.execute(
        "SELECT distance FROM stop_distance WHERE route_id=%s AND direction_id=%s AND from_stop_id=%s AND to_stop_id=%s",
        [routeID, directionID, fromStop, toStop]
    )
    row = cur.fetchone()
    conn.commit()
    cur.close()
    if len(row) == 0:
        cur = conn.cursor()
        cur.execute(
            "SELECT AVG(distance) FROM stop_distance WHERE route_id=%s GROUP BY route_id;",
            [routeID]
        )
        row = cur.fetchone()
        conn.commit()
        cur.close()
        if len(row) == 0:
            conn.close()
            return AVG_DISTANCE
    conn.close()
    return row[0]
