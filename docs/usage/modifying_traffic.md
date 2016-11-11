# Modifying traffic

In *modify mode*, Hoverfly does not make use of its cache. Traffic is passed through Hoverfly, and middleware is executed both on out-going and in-bound traffic. In this mode, Middleware can be used to manipulate any part of a request or a response.

This is a mode with potential applications outside of simulating external services for development and testing. Possible use-cases include introducing transparent redirects and injecting authentication headers.

## Set up

You will need to have your application or OS configured to use Hoverfly as a proxy (see the **Installation and setup** section).

To modify traffic on the fly, you do not need to have captured any traffic or imported any Hoverfly JSON. However, it is important that you read the **Middleware** section first, if you haven't already.

For the example below, you will also need Python installed on the Hoverfly host.

## Modify some traffic

### Create a middleware script
First, create a middleware script to modify traffic. Save the following code in a file named `modify.py` in your Hoverfly directory and make it executable with `chmod +x modify.py`:

    #!/usr/bin/env python

    import sys
    import json
    import logging

    logging.basicConfig(filename='middleware.log', level=logging.DEBUG)
    logging.debug('Middleware "modify_request" called')


    def main():
        data = sys.stdin.readlines()
        # this is a json string in one line so we are interested in that one line
        payload = data[0]
        logging.debug(payload)

        payload_dict = json.loads(payload)

        payload_dict['request']['destination'] = "mirage.readthedocs.org"
        payload_dict['request']['method'] = "GET"

        payload_dict['response']['status'] = 201
        # returning new payload
        print(json.dumps(payload_dict))

    if __name__ == "__main__":
        main()

This middleware script will transparently redirect **any** request to **any host** that is passed through Hoverfly to `http://mirage.readthedocs.org`.

### Put Hoverfly into Modify mode

There are three ways of doing this:

1. Start Hoverfly with the `-modify` flag and specify middleware with the `-middleware` flag:

       ./hoverfly -modify -middleware "modify.py"

2. Ensure Hoverfly is running (in any mode), and that the middleware has been specified, then select "modify" in the Admin UI which is available at `http://${HOVERFLY_HOST}:8888`.

3. Ensure Hoverfly is running (in any mode), and that the middleware has been specified, then make an API call:

       curl -H "Content-Type application/json" -X POST -d '{"mode":"modify"}' http://${HOVERFLY_HOST}:8888/api/v2/hoverfly/mode

Now cURL any (HTTP) URL with Hoverfly as a proxy. For example:

    curl http://flask.readthedocs.org/ --proxy http://${HOVERFLY_HOST}:8500/

You will see that the response comes from `http://mirage.readthedocs.org`.

## What is happening?

Hoverfly is executing the `modify.py` script on the request and the response. The middleware takes the service data JSON string via STDIN, and replaces the `destination` with the string `mirage.readthedocs.org`. It also takes the response JSON string via STDIN and replaces the `method` with `201`.

## Further reading

This example is taken from a more detailed step-by-step guide:

[Modifying traffic on the fly](https://specto.io/blog/service-virtualization-is-so-last-year.html)
