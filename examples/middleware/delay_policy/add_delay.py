#!/usr/bin/env python
import sys
import logging
from time import sleep

logging.basicConfig(filename='delay_middleware.log', level=logging.DEBUG)
logging.debug('Delay middleware is called')

SLEEP_SECS = 5

def main():

    data = sys.stdin.readlines()
    # this is a json string in one line so we are interested in that one line
    payload = data[0]
    logging.debug("sleeping for %s seconds" % SLEEP_SECS)
    sleep(SLEEP_SECS)


    # do not modifying payload, returning same one
    print(payload)

if __name__ == "__main__":
    main()


