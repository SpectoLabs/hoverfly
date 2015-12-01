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
 
## License

Apache License version 2.0 [See LICENSE for details](https://github.com/SpectoLabs/hoverfly/blob/master/LICENSE).

(c) [SpectoLabs](https://specto.io) 2015.
