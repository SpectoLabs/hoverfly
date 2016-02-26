#!/usr/bin/env python
import sys
import logging
import json


logging.basicConfig(filename='middleware_flight_search.log', level=logging.DEBUG)


def get_trip_info(origin, destination, date):
    """
    Provides basic template for response, you can change as many things as you like.
    :param origin: from which airport your trip beings
    :param destination: where are you flying to
    :param date: when
    :return:
    """
    template = {
        "kind": "qpxExpress#tripsSearch",
        "trips": {
            "kind": "qpxexpress#tripOptions",
            "requestId": "SYzLMFMFPCrebUp5H0NaGL",
            "data": {
                "kind": "qpxexpress#data",
                "airport": [
                    {
                        "kind": "qpxexpress#airportData",
                        "code": "AMS",
                        "city": "AMS",
                        "name": "Amsterdam Schiphol Airport"
                    },
                    {
                        "kind": "qpxexpress#airportData",
                        "code": "LGW",
                        "city": "LON",
                        "name": "London Gatwick"
                    }
                ],
                "city": [
                    {
                        "kind": "qpxexpress#cityData",
                        "code": "AMS",
                        "name": "Amsterdam"
                    },
                    {
                        "kind": "qpxexpress#cityData",
                        "code": "LON",
                        "name": "London"
                    }
                ],
                "aircraft": [
                    {
                        "kind": "qpxexpress#aircraftData",
                        "code": "319",
                        "name": "Airbus A319"
                    },
                    {
                        "kind": "qpxexpress#aircraftData",
                        "code": "320",
                        "name": "Airbus A320"
                    }
                ],
                "tax": [
                    {
                        "kind": "qpxexpress#taxData",
                        "id": "GB_001",
                        "name": "United Kingdom Air Passengers Duty"
                    },
                    {
                        "kind": "qpxexpress#taxData",
                        "id": "UB",
                        "name": "United Kingdom Passenger Service Charge"
                    }
                ],
                "carrier": [
                    {
                        "kind": "qpxexpress#carrierData",
                        "code": "BA",
                        "name": "British Airways p.l.c."
                    }
                ]
            },
            "tripOption": [
                {
                    "kind": "qpxexpress#tripOption",
                    "saleTotal": "GBP47.27",
                    "id": "OAcAQw8rr9MNhwQoBntUKJ001",
                    "slice": [
                        {
                            "kind": "qpxexpress#sliceInfo",
                            "duration": 75,
                            "segment": [
                                {
                                    "kind": "qpxexpress#segmentInfo",
                                    "duration": 75,
                                    "flight": {
                                        "carrier": "BA",
                                        "number": "2762"
                                    },
                                    "id": "GStLakphRYJX3LbK",
                                    "cabin": "COACH",
                                    "bookingCode": "O",
                                    "bookingCodeCount": 1,
                                    "marriedSegmentGroup": "0",
                                    "leg": [
                                        {
                                            "kind": "qpxexpress#legInfo",
                                            "id": "LgJHYCVgG0AiE1PH",
                                            "aircraft": "320",
                                            "arrivalTime": "%sT18:05+01:00" % date,
                                            "departureTime": "%sT15:50+00:00" % date,
                                            "origin": origin,
                                            "destination": destination,
                                            "originTerminal": "N",
                                            "duration": 75,
                                            "mileage": 226,
                                            "meal": "Snack or Brunch"
                                        }
                                    ]
                                }
                            ]
                        }
                    ],
                    "pricing": [
                        {
                            "kind": "qpxexpress#pricingInfo",
                            "fare": [
                                {
                                    "kind": "qpxexpress#fareInfo",
                                    "id": "A855zsItBCELBykaeqeBDQb5hPZQIOtkOZ8uDq0lD5VU",
                                    "carrier": "BA",
                                    "origin": origin,
                                    "destination": destination,
                                    "basisCode": "OV1KO"
                                }
                            ],
                            "segmentPricing": [
                                {
                                    "kind": "qpxexpress#segmentPricing",
                                    "fareId": "A855zsItBCELBykaeqeBDQb5hPZQIOtkOZ8uDq0lD5VU",
                                    "segmentId": "GStLakphRYJX3LbK"
                                }
                            ],
                            "baseFareTotal": "GBP22.00",
                            "saleFareTotal": "GBP22.00",
                            "saleTaxTotal": "GBP25.27",
                            "saleTotal": "GBP47.27",
                            "passengers": {
                                "kind": "qpxexpress#passengerCounts",
                                "adultCount": 1
                            },
                            "tax": [
                                {
                                    "kind": "qpxexpress#taxInfo",
                                    "id": "UB",
                                    "chargeType": "GOVERNMENT",
                                    "code": "UB",
                                    "country": "GB",
                                    "salePrice": "GBP12.27"
                                },
                                {
                                    "kind": "qpxexpress#taxInfo",
                                    "id": "GB_001",
                                    "chargeType": "GOVERNMENT",
                                    "code": "GB",
                                    "country": "GB",
                                    "salePrice": "GBP13.00"
                                }
                            ],
                            "fareCalculation": "LON BA AMS 33.71OV1KO NUC 33.71 END ROE 0.652504 FARE GBP 22.00 XT 13.00GB 12.27UB",
                            "latestTicketingTime": "2016-01-11T23:59-05:00",
                            "ptc": "ADT"
                        }
                    ]
                },
                {
                    "kind": "qpxexpress#tripOption",
                    "saleTotal": "GBP62.27",
                    "id": "OAcAQw8rr9MNhwQoBntUKJ002",
                    "slice": [
                        {
                            "kind": "qpxexpress#sliceInfo",
                            "duration": 80,
                            "segment": [
                                {
                                    "kind": "qpxexpress#segmentInfo",
                                    "duration": 80,
                                    "flight": {
                                        "carrier": "BA",
                                        "number": "2758"
                                    },
                                    "id": "GW8rUjsDA234DdHV",
                                    "cabin": "COACH",
                                    "bookingCode": "Q",
                                    "bookingCodeCount": 9,
                                    "marriedSegmentGroup": "0",
                                    "leg": [
                                        {
                                            "kind": "qpxexpress#legInfo",
                                            "id": "Lp08eKxnXnyWfJo4",
                                            "aircraft": "319",
                                            "arrivalTime": "%sT10:05+01:00" % date,
                                            "departureTime": "%sT07:45+00:00" % date,
                                            "origin": origin,
                                            "destination": destination,
                                            "originTerminal": "N",
                                            "duration": 80,
                                            "mileage": 226,
                                            "meal": "Snack or Brunch"
                                        }
                                    ]
                                }
                            ]
                        }
                    ],
                    "pricing": [
                        {
                            "kind": "qpxexpress#pricingInfo",
                            "fare": [
                                {
                                    "kind": "qpxexpress#fareInfo",
                                    "id": "AslXz8S1h3mMcnYUQ/v0Zt0p9Es2hj8U0We0xFAU1qDE",
                                    "carrier": "BA",
                                    "origin": origin,
                                    "destination": destination,
                                    "basisCode": "QV1KO"
                                }
                            ],
                            "segmentPricing": [
                                {
                                    "kind": "qpxexpress#segmentPricing",
                                    "fareId": "AslXz8S1h3mMcnYUQ/v0Zt0p9Es2hj8U0We0xFAU1qDE",
                                    "segmentId": "GW8rUjsDA234DdHV"
                                }
                            ],
                            "baseFareTotal": "GBP37.00",
                            "saleFareTotal": "GBP37.00",
                            "saleTaxTotal": "GBP25.27",
                            "saleTotal": "GBP62.27",
                            "passengers": {
                                "kind": "qpxexpress#passengerCounts",
                                "adultCount": 1
                            },
                            "tax": [
                                {
                                    "kind": "qpxexpress#taxInfo",
                                    "id": "UB",
                                    "chargeType": "GOVERNMENT",
                                    "code": "UB",
                                    "country": "GB",
                                    "salePrice": "GBP12.27"
                                },
                                {
                                    "kind": "qpxexpress#taxInfo",
                                    "id": "GB_001",
                                    "chargeType": "GOVERNMENT",
                                    "code": "GB",
                                    "country": "GB",
                                    "salePrice": "GBP13.00"
                                }
                            ],
                            "fareCalculation": "LON BA AMS 56.70QV1KO NUC 56.70 END ROE 0.652504 FARE GBP 37.00 XT 13.00GB 12.27UB",
                            "latestTicketingTime": "%sT23:59-05:00" % date,
                            "ptc": "ADT"
                        }
                    ]
                }
            ]
        }
    }

    return template


def get_origin(body):
    """
    Getting origin
    :param body: request payload as dict
    :return:
    """
    return body['request']['slice'][0]['origin']


def get_destination(body):
    """
    Getting destination
    :param body: request payload as dict
    :return:
    """
    return body['request']['slice'][0]['destination']


def get_date(body):
    return body['request']['slice'][0]['date']


def main():
    data = sys.stdin.readlines()
    # this is a json string in one line so we are interested in that one line
    payload = data[0]
    # logging.debug(payload)

    payload_dict = json.loads(payload)

    # getting requested city and date
    request_body = json.loads(payload_dict['request']['body'])
    trip_origin = get_origin(request_body)
    trip_destination = get_destination(request_body)
    trip_date = get_date(request_body)

    # some debugging info
    logging.debug("Getting trip! Origin: %s, destination: %s. Date: %s" % (trip_origin, trip_destination, trip_date))
    logging.debug(payload)

    # adding json header
    payload_dict['response']['headers '] = {"Content-Type": ["application/json"]}

    if len(trip_origin) > 3 or len(trip_destination) > 3:
        # checking whether origin and destination doesn't exceed 3 symbols
        payload_dict['response']['status'] = 400
        payload_dict['response']['body'] = '{"error": "origin and destination cannot exceed 3 symbols limit"}'
        print(json.dumps(payload_dict))
        return

    # preparing response
    payload_dict['response']['status'] = 200
    payload_dict['response']['body'] = json.dumps(get_trip_info(
        origin=trip_origin,
        destination=trip_destination,
        date=trip_date))

    # returning new payload
    print(json.dumps(payload_dict))

if __name__ == "__main__":
    main()


