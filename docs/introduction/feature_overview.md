# Feature overview

* "Capture" traffic between a client and a server application
* Use captured traffic to simulate the server application
* Export captured service data as a JSON file
* Import service data JSON files
* Simulate latency by specifying delays which can be applied to individual URLs based on regex patterns, or based on HTTP method
* Flexible request matching using templates 
* Supports "middleware" (which can be written in any language) to manipulate data in requests or responses, or to simulate unexpected behaviour such as malformed responses or random errors
* Supports local or remote middleware execution (for example on [AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/welcome.html))
* Uses [BoltDB](https://github.com/boltdb/bolt) to persist data in a binary file on disk - so no additional database is required
* REST API
* Run as a transparent proxy or as a webserver 
* High performance with minimal overhead
* JUnit rule "wrapper" is available as a Maven dependency
* Supports HTTPS and can generate certificates if required
* Authentication (combination of Basic Auth and [JWT](https://jwt.io/))
* Command line interface ("hoverctl")
* Admin UI to change state and view basic metrics

