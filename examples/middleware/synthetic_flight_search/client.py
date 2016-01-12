import requests
import json
from argparse import ArgumentParser
from pprint import pprint as pp

url = "https://www.googleapis.com/qpxExpress/v1/trips/search?key=XXX"


def get_payload(origin, destination, date):
    payload = {
        "request": {
            "slice": [
                {
                    "origin": origin,
                    "destination": destination,
                    "date": date,
                    "maxStops": 0
                }
            ],
            "passengers": {
                "adultCount": 1,
                "infantInLapCount": 0,
                "infantInSeatCount": 0,
                "childCount": 0,
                "seniorCount": 0
            },
            "solutions": 2,
            "refundable": "false"
        }
    }
    return payload


def main():
    parser = ArgumentParser(description="Perform query to https://qpx-express-demo.itasoftware.com/")
    parser.add_argument("--origin", help="origin airport 'LGW' ")
    parser.add_argument("--destination", help="destination airport 'AMS'")
    parser.add_argument("--date", help="date of travel '2016-06-20'")

    args = parser.parse_args()

    response = requests.post(url, json.dumps(get_payload(origin=args.origin,
                                                         destination=args.destination,
                                                         date=args.date)),
                             headers={'content-type': 'application/json'}, verify=False)
    pp(response.json())


if __name__ == "__main__":
    main()
