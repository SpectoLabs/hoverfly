#!/usr/bin/env python

import requests
from argparse import ArgumentParser
import pickle
 
limit = 50

# getting urls and dumping them into file
def get_urls():
    sites = requests.get("http://readthedocs.org/api/v1/project/?limit=%s&amp;offset=0&amp;format=json" % limit)
 
    objects = sites.json()['objects']
 
    links = ["http://readthedocs.org" + x['resource_uri'] for x in objects]
 
    with open("links.p", "wb") as outfile:
        pickle.dump(links, outfile)
 
 
def fetch_links():
    with open("links.p", "rb") as infile:
        links = pickle.load(infile)
 
    import time
    start = time.time()
 
    for link in links:
        response = requests.get(link)
        print("url: %s, status code: %s" % (link, response.status_code))
 
    print(time.time() - start)
 
 
# main function
def main():
    parser = ArgumentParser(description="Perform proxy testing/URL list creation")
    parser.add_argument("--urls", help="download and save urls ")
    args = parser.parse_args()
 
    # get urls
    if args.urls:
        get_urls()
    else:
        fetch_links()
 
 
if __name__ == "__main__":
    main()
