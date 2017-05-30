#!/usr/bin/env python
import sys
import logging
import json
from time import gmtime, strftime


logging.basicConfig(filename='middleware_synthetic.log', level=logging.DEBUG)
logging.debug('Middleware "synthetic service" called')


def main():
    data = sys.stdin.readlines()
    # this is a json string in one line so we are interested in that one line
    payload = data[0]
    logging.debug(payload)

    payload_dict = json.loads(payload)

    dest = payload_dict['request']['destination']

    payload_dict['response']['status'] = 200
    payload_dict['response']['body'] = "You called (%s). I am synthethic service, maybe I could do more?\n" \
                                       "Current time: %s"

    # returning new payload
    print(json.dumps(payload_dict))

if __name__ == "__main__":
    main()


