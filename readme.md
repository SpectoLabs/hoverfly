# Hoverfly: dependencies without the sting

Hoverfly is an experiment in lightweight, open source [service virtualization](https://en.wikipedia.org/wiki/Service_virtualization). Using Hoverfly, you can virtualize your application dependencies to create a self-contained development/test environment.

Hoverfly is a proxy written in Go. It can capture HTTP(s) traffic between an application under test and external services, and then replace the external services. Another powerful feature - middleware modules, where users can introduce their own custom logic. Middleware modules can be written in any language. Hoverfly uses Redis for persistence.

## Installation

You can just grab a binary from releases or you can build it yourself. Use [glide](https://github.com/Masterminds/glide) to
fetch dependencies with:
* glide up

Then, build it:
* go build

Run it:
* ./hoverfly


## Destination configuration

Specifying which site to record/virtualize with regular expression (by default it processes everything):

    ./hoverfly --destination="."

## Modes (Capture / Virtualize / Synthesize)

Hoverfly has different operating modes. Each mode changes proxy behavior. Based on current mode Hoverfly selects
whether it should capture requests/response, look for them in the cache or send them directly into middleware and respond with a payload that middleware generated (more on middleware below).

### Virtualize

By default proxy starts in virtualize mode. You can apply middleware to each response.

### Capture

When capture mode is active - Hoverfly acts as a man in the middle. It makes requests on behalf of a client and records responses, response is then sent back to the original client.

To switch to capture mode, add "--capture" flag during startup:

    ./hoverfly --capture

Or you can use API call to change proxy state while running.


Do a curl request with proxy details:

    curl http://mirage.readthedocs.org --proxy http://localhost:8500/

###  Synthesize

Hoverfly can create responses to requests on the fly. Synthesize mode intercepts requests (also respects --destination flag)
and applies supplied middleware (user is required to supply --middleware flag, you can read more about it below). Middleware
is expected to populate response payload, so Hoverfly can then reconstruct it and return it to the client. Example of synthetic
service can be found in _this_repo/examples/middleware/synthetic_service/synthetic.py_. You can test it out by running:

    ./hoverfly --synthesize --middleware "./examples/middleware/synthetic_service/synthetic.py"

### Modify

This mode is not yet implemented, it will be able to modify outgoing requests and incoming responses.

## HTTPS record

Add ca.pem to your trusted certificates or turn off verification, with curl you can make insecure requests with -k:

    curl https://www.bbc.co.uk --proxy http://localhost:8500 -k

## API

Access administrator API under default port 8888:

* Recorded requests: GET http://proxy_hostname:8888/records ( __curl http://proxy_hostname:8888/records__ )
* Wipe cache: DELETE http://proxy_hostname:8888/records ( __curl -X DELETE http://proxy_hostname:8888/records__ )
* Get current proxy state: GET http://proxy_hostname:8888/state ( __curl http://proxy_hostname:8888/state__ )
* Set proxy state: POST http://proxy_hostname:8888/state, where
   + body to start virtualizing: {"mode":"virtualize"}
   + body to start capturing: {"mode":"capture"}
* Exporting recorded requests to a file: __curl http://proxy_hostname:8888/records > requests.json__
* Importing requests from file: __curl --data "@/path/to/requests.json" http://localhost:8888/records__


## Middleware

Hoverfly supports (experimental feature) external middleware modules. You can write them in __any language you want__ !.
These middleware modules are expected to take standard input (stdin) and should output to stdout same structured JSON string.
Payload example:

```javascript
{
	"response": {
		"status": 200,
		"body": "body here",
		"headers": {
			"Content-Type": ["text/html"],
			"Date": ["Tue, 01 Dec 2015 16:49:08 GMT"],
		}
	},
	"request": {
		"path": "/",
		"method": "GET",
		"destination": "1stalphaomega.readthedocs.org",
		"query": ""
	},
	"id": "5d4f6b1d9f7c2407f78e4d4e211ec769"
}
```
Middleware is executed only when request is matched so for fully dynamic responses where you are
generating response on the fly you can use ""--synthesize" flag.

In order to use your middleware, just add path to executable:

    ./hoverfly --middleware "./examples/middleware/modify_response/modify_response.py"

Basic example of a Python module to change response body and add 2 second delay:

```python
#!/usr/bin/env python
import sys
import logging
import json
from time import sleep


logging.basicConfig(filename='middleware.log', level=logging.DEBUG)
logging.debug('Middleware is called')


def main():
    data = sys.stdin.readlines()
    # this is a json string in one line so we are interested in that one line
    payload = data[0]
    logging.debug(payload)

    payload_dict = json.loads(payload)

    payload_dict['response']['status'] = 201
    payload_dict['response']['body'] = "body was replaced by middleware"

    # now let' sleep for 2 seconds
    sleep(2)

    # returning new payload
    print(json.dumps(payload_dict))

if __name__ == "__main__":
    main()

```

Save this file with python extension, _chmod +x_ it and run hoverfly:

    ./hoverfly --middleware "./this_file.py"
    

### How middleware interacts with different modes?

Each mode is affected by middleware in different ways. Since JSON payload has request and response structures - some middleware
 will not change anything, some information about middleware behaviour:
  * __Capture Mode__: middleware affects only outgoing requests. 
  * __Virtualize Mode__: middleware affects only responses (cache contents remain untouched).
  * __Synthesize Mode__: middleware creates responses.
  * __Modify Mode__: middleware affects requests and responses.
  


## Debugging

You can supply "-v" flag to enable verbose logging.


## License

Apache License version 2.0 [See LICENSE for details](https://github.com/SpectoLabs/hoverfly/blob/master/LICENSE).

(c) [SpectoLabs](https://specto.io) 2015.
