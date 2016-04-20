# Feature Overview

Hoverfly is an HTTP(s) proxy written in Go. It is available as a single binary executable file. An overview of the main features:

* "Capture" real traffic between a client and a server application
* Use captured traffic to simulate the server application
* Export captured traffic as a JSON file
* Import service definition (JSON) files which can be created manually or can have been exported previously
* Supports "middleware" (which can be written in any language) to manipulate data in requests or responses, or to simulate network failure and latency etc.
* Uses BoltDB to persist data in a binary file on disk - so no additional database is required
* REST API
* Highly performant and uses very little resource
* JUnit rule "wrapper" is available as a Maven dependency
* Supports HTTPS and can generate certificates if required
* Admin UI (with authentication) to change state and view basic metrics
