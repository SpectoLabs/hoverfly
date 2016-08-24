#!/usr/bin/env python
import sys
import logging
import json


logging.basicConfig(filename='middleware.log', level=logging.DEBUG)
logging.debug('Middleware is called')


def main():
    data = sys.stdin.readlines()
    # this is a json string in one line so we are interested in that one line
    payload = data[0]
    logging.debug(payload)

    payload_dict = json.loads(payload)

    payload_dict['response']['status'] = 201
    payload_dict['response']['body'] = "body was replaced by middleware\n"

    # returning new payload
    print(json.dumps(payload_dict))

if __name__ == "__main__":
    main()


