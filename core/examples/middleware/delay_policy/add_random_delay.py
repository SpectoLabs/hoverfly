#!/usr/bin/env python
import sys
import logging
import random
from time import sleep

logging.basicConfig(filename='random_delay_middleware.log', level=logging.DEBUG)
logging.debug('Random delay middleware is called')

# set delay to random value less than one second

SLEEP_SECS = random.random()

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
