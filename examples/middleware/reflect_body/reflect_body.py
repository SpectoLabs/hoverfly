#!/usr/bin/env python
import sys
import json


def main():
    """
    Simple middleware to reflect back request body to response.
    """
    data = sys.stdin.readlines()
    # this is a json string in one line so we are interested in that one line
    payload = data[0]

    payload_dict = json.loads(payload)

    payload_dict['response']['body'] = payload_dict['request']['body']
    payload_dict['response']['status'] = 200

    # returning new payload
    print(json.dumps(payload_dict))

if __name__ == "__main__":
    main()
