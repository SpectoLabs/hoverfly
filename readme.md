# Hoverfly: dependencies without the sting

Hoverfly is an experiment in lightweight, open source [service virtualization](https://en.wikipedia.org/wiki/Service_virtualization). Using Hoverfly, you can virtualize your application dependencies to create a self-contained development/test environment. 

Hoverfly is a transparent proxy written in Go. It can capture HTTP(s) traffic between an application under test and external services, and then replace the external services. Hoverfly uses Redis for persistence.


## Configuration

Specifying which site to capture/virtualize with regular expression (by default it captures everything):

    ./hoverfly --destination="."

By default proxy is always in virtualize mode. To switch to capture mode, add "--capture" flag during startup:

    ./hoverfly --capture

Or you can use API call to change proxy state while running.


Do a curl request with proxy details: 

    curl http://mirage.readthedocs.org --proxy http://localhost:8500/

### HTTPS capture

Add ca.pem to your trusted certificates or turn off verification, with curl you can make insecure requests with -k: 

    curl https://www.bbc.co.uk --proxy http://localhost:8500 -k

### Virtualizing services

Start proxy in virtualize mode:

    ./hoverfly

## API

Access administrator API under default port 8888:

* Recorded requests: GET http://proxy_hostname:8888/records ( __curl http://proxy_hostname:8888/records__ )
* Wipe cache: DELETE http://proxy_hostname:8888/records ( __curl -X DELETE http://proxy_hostname:8888/records__ ) 
* Get current proxy state: GET http://proxy_hostname:8888/state ( __curl http://proxy_hostname:8888/state__ )
* Set proxy state: POST http://proxy_hostname:8888/state, where
   + body to start playback: {"record":false}
   + body to start recording: {"record":true}
* Exporting recorded requests to a file: __curl http://proxy_hostname:8888/records > requests.json__
* Importing requests from file: __curl --data "@/path/to/requests.json" http://localhost:8888/records__


## Middleware

Hoverfly supports (experimental feature) external middleware modules. You can write them in __any language you want!__.
Middleware is executed when request is matched and response found in cache so for fully dynamic responses where you are 
generating response on the fly - just add dummy request through import functionality. 

In order to use your middleware, just add path to executable: 
* ./hoverfly --middleware "./examples/middleware/modify_response/modify_response.py" 

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
*./hoverfly --middleware "./this_file.py"


 
## License

Apache License version 2.0 [See LICENSE for details](https://github.com/SpectoLabs/hoverfly/blob/master/LICENSE).

(c) [SpectoLabs](https://specto.io) 2015.
