from lib.args import extract_route_id


def main():
    # Get route_id from CLI arguments
    route_id = extract_route_id()
    # Print route_id as test
    print(route_id)
