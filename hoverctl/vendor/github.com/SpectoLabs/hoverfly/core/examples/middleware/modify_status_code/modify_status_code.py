#!/usr/bin/env python
import sys
import logging
import json


logging.basicConfig(filename='middleware_status_code.log', level=logging.DEBUG)
logging.debug('Middleware to change status code called')


def main():
    data = sys.stdin.readlines()
    # this is a json string in one line so we are interested in that one line
    payload = data[0]
    logging.debug(payload)

    payload_dict = json.loads(payload)

    payload_dict['response']['status'] = 301

    # returning new payload
    print(json.dumps(payload_dict))

if __name__ == "__main__":
    main()


