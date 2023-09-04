#!/usr/bin/env python
import sys
import requests


def get_ip_info():
    url = "http://ip.jsontest.com/"

    try:
        response = requests.get(url)
        response.raise_for_status()  # Raise an exception if the request was not successful (HTTP status code >= 400)

        data = response.json()  # Parse the JSON response
        return data
    except requests.exceptions.RequestException as e:
        print(f"An error occurred: {e}")
        return None

def main():

    data = sys.stdin.readlines()
    # this is a request-response pair json string. We can use this pair to make outbound call or do any operations
    payload = data[0]
    ip_info = get_ip_info()
    if ip_info:
            print(f"HTTP call invoked from IP Address: {ip_info['ip']}")
    else:
            print("Failed to retrieve IP information.")

if __name__ == "__main__":
    main()