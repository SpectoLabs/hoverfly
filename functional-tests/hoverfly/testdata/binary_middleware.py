#!/usr/bin/env python
import sys
import json
import logging
import base64
import os

logging.basicConfig(filename='middleware_request.log', level=logging.DEBUG)
logging.debug('Middleware "modify_request" called')


def main():
    data = sys.stdin.readlines()
    # this is a json string in one line so we are interested in that one line
    payload = data[0]
    logging.debug(payload)

    payload_dict = json.loads(payload)

    bytes = open(os.getcwd() + "/testdata/1x1.png").read()
    responseBody = base64.b64encode(bytes) # We need to base 64 encode binary data
    contentLength = len(bytes) # Content length should be what we expect the unencoded binary data to be

    payload_dict['request']['body'] = "CHANGED_REQUEST_BODY"

    payload_dict['response']['status'] = 200
    payload_dict['response']['headers'] = {'Content-Length': [str(contentLength)], 'Content-Type' : ["image/png"]}
    payload_dict['response']['body'] = responseBody
    payload_dict['response']['encodedBody'] = True

    # returning new payload
    print(json.dumps(payload_dict))

if __name__ == "__main__":
    main()

