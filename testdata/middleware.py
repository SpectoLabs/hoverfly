#!/usr/bin/env python
import sys
import json
import logging

logging.basicConfig(filename='middleware_request.log', level=logging.DEBUG)
logging.debug('Middleware "modify_request" called')


def main():
    data = sys.stdin.readlines()
    # this is a json string in one line so we are interested in that one line
    payload = data[0]
    logging.debug(payload)

    payload_dict = json.loads(payload)

    payload_dict['response']['body'] = "CHANGED_RESPONSE_BODY"
    payload_dict['request']['body'] = "CHANGED_REQUEST_BODY"
    payload_dict['response']['status'] = 200
    payload_dict['response']['headers'] = {'Content-Length': ["21"]}

    # returning new payload
    print(json.dumps(payload_dict))

if __name__ == "__main__":
    main()

