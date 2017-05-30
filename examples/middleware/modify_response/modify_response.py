#!/usr/bin/env python
import sys
import json

def main():
    data = sys.stdin.readlines()
    payload = data[0]

    payload_dict = json.loads(payload)
    payload_dict['response']['body'] = "body was replaced by middleware\n"

    print(json.dumps(payload_dict))

if __name__ == "__main__":
    main()
