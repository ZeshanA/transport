from datetime import timedelta

from db.stop_distance import get_stop_distance
from lib.date import parse_datetime, format_datetime
from lib.model import predict


def calculate(model, request):
    journey = request['journey']
    route_id, direction_id = journey['routeID'], journey['directionID']
    stop_list, from_stop, to_stop = request['stopList'], journey['fromStop'], journey['toStop']
    route_segment = get_route_segment(stop_list, from_stop, to_stop)
    arrival_time = parse_datetime(journey['arrivalTime'])
    avg_stop_to_stop_duration = int(request['averageJourneyTime'] / len(route_segment))
    estimated_start_time = arrival_time - timedelta(seconds=request['averageJourneyTime'])
    journey_durations, arrival_stamps = [], [estimated_start_time]
    # Loop over the possible "starting" points
    for i in range(0, len(route_segment) - 1):
        source_stop, dest_stop = route_segment[i], route_segment[i + 1]
        mvmt = dict(request['sampleMovement'])
        mvmt['Latitude'], mvmt['Longitude'] = source_stop['latitude'], source_stop['longitude']
        mvmt['DistanceFromStop'] = get_stop_distance(route_id, direction_id, source_stop['id'], dest_stop['id'])
        mvmt['ExpectedArrivalTime'] = format_datetime(
            arrival_stamps[-1] + timedelta(seconds=avg_stop_to_stop_duration))
        mvmt['StopPointRef'] = dest_stop['id']
        if i == 0:
            mvmt['Timestamp'] = format_datetime(estimated_start_time)
        else:
            mvmt['Timestamp'] = format_datetime(arrival_stamps[-1] + timedelta(seconds=journey_durations[-1]))
        duration = predict(model, dict(mvmt))
        journey_durations.append(duration)
        arrival_stamps.append(arrival_stamps[-1] + timedelta(seconds=duration))
    return sum(journey_durations)


def get_route_segment(stop_list, from_stop, to_stop):
    from_stop_idx, to_stop_idx = 0, 0
    for i, stop in enumerate(stop_list):
        if stop['id'] == from_stop:
            from_stop_idx = i
        if stop['id'] == to_stop:
            to_stop_idx = i
    return stop_list[from_stop_idx:to_stop_idx + 1]
