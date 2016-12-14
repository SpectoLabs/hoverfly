Installation and setup
----------------------

Hoverfly is a single binary file. It comes with an optional command line
interface tool called hoverctl.

Download one of the zip files below, extract the Hoverfly and hoverctl
binaries, and move them to a directory on your
`PATH <https://www.java.com/en/download/help/path.xml>`__.

-  `macOS
   64bit <https://github.com/SpectoLabs/hoverfly/releases/download/v0.9.0/hoverfly_bundle_OSX_amd64.zip>`__
-  `Linux
   32bit <https://github.com/SpectoLabs/hoverfly/releases/download/v0.9.0/hoverfly_bundle_linux_386.zip>`__
-  `Linux
   64bit <https://github.com/SpectoLabs/hoverfly/releases/download/v0.9.0/hoverfly_bundle_linux_amd64.zip>`__
-  `Windows
   32bit <https://github.com/SpectoLabs/hoverfly/releases/download/v0.9.0/hoverfly_bundle_windows_386.zip>`__
-  `Windows
   64bit <https://github.com/SpectoLabs/hoverfly/releases/download/v0.9.0/hoverfly_bundle_windows_amd64.zip>`__

Use `Homebrew <http://brew.sh/>`__ to install Hoverfly and hoverctl:

::

    brew install SpectoLabs/tap/hoverfly

Docker image
~~~~~~~~~~~~

The Hoverfly docker image contains only the Hoverfly binary.

::

    docker run -d -p 8888:8888 -p 8500:8500 spectolabs/hoverfly

The port mapping is for the Hoverfly AdminUI/API and proxy respectively.
The Docker image supports Hoverfly flags (see **Flags and environment
variables** in the **Usage** section).

JUnit rule
~~~~~~~~~~

The Hoverfly JUnit rule is available as a Maven dependency. To get
started, add the following to your pom.xml:

::

    <groupId>io.specto</groupId>
    <artifactId>hoverfly-junit</artifactId>
    <version>0.1.4</version>

This will download the JUnit rule and the Hoverfly binary. More
information on how to use the Hoverfly JUnit rule is available here:

`Easy API Simulation with the Hoverfly JUnit
Rule <https://specto.io/blog/hoverfly-junit-api-simulation.html>`__

Setting the HOVERFLY\_HOST environment variable
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Throughout the documentation, ``${HOVERFLY_HOST}`` is used in API
examples and Admin UI guides. If you are running a binary on your local
machine, the Hoverfly host will be ``localhost``. However, if you are
running Hoverfly on a remote machine, or using Docker Machine for
example, it will be different.

To make things easier when following the documentation, it is
recommended that you set the HOVERFLY\_HOST environment variable. For
example:

::

    export HOVERFLY_HOST=localhost

Hoverfly as an HTTP(S) proxy
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Hoverfly is primarily a proxy - although it can run (with limitations)
as a webserver. To use Hoverfly in your development or test environment
to capture traffic, you need to ensure that your application is using it
as a proxy. This can be done at the OS level, or at the application
level, depending on the environment.

-  `Windows proxy settings explained <http://blog.raido.be/?p=426>`__
-  `Firefox proxy
   setting <https://support.mozilla.org/en-US/kb/advanced-panel-settings-in-firefox#w_connection>`__
-  `Java Networking and
   Proxies <https://docs.oracle.com/javase/6/docs/technotes/guides/net/proxies.html>`__

Admin UI
~~~~~~~~

The Hoverfly Admin UI is available on port 8888 by default:

::

    http://${HOVERFLY_HOST}:8888

The port is configurable (see the **Flags and environment variables**
section).

When authentication is disabled (which is the default), you can use
**any username and password combination** to access the Admin UI.

The Admin UI can be used to change the Hoverfly mode, and to view basic
analytic information. It uses the Hoverfly API.

Starting Hoverfly as a webserver
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

While Hoverfly is primarily a proxy, there are some situations in which
you may need to run it as a webserver - for example if setting the host
OS or application to use a proxy is not possible or desirable.

Currently, when running as a webserver, Hoverfly can only be set to
*simulate* mode. This is useful if you have a simulation that has been
created by capturing traffic using a Hoverfly instance running as a
proxy, and you want to import and run it in an environment which cannot
be configured to use a proxy.

If you are using hoverctl to manage your instance of Hoverfly, you can
start Hoverfly as a webserver using hoverctl.

::

    hoverctl start webserver

If you are running the Hoverfly binary, you can specify the webserver
flag which will start Hoverfly as a webserver.

::

    ./hoverfly -webserver

**NOTE:** Currently HTTPS is not supported when running Hoverfly as a
webserver. HTTPS support when running as a webserver is on the roadmap.

Simulations and request matching
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

When running as a webserver, although Hoverfly functionality is limited
to simulate mode, Hoverfly still uses the standard simulations.

When Hoverfly is running as a webserver, the server is available on the
same port as the proxy. When you make requests to the webserver,
responses will be matched in the same way as with the proxy, except the
host will be disregarded. This means that if you have a simulation with
multiple hosts, they will all be served from the same host.

--------------

Capturing traffic
-----------------

Set up
~~~~~~

To capture traffic, you will need:

-  A client application (this could be an application you are developing
   or testing, or cURL if you just want to experiment)
-  An external service or API that your client application communicates
   with over HTTP or HTTPS
-  Hoverfly - either as a binary or the Docker image
-  Your application or OS configured to use Hoverfly as an HTTP or HTTPS
   proxy

Please follow the steps in the **Installation and setup** section if you
haven't already.

If you want to capture traffic to and from a service over HTTPS, you
will need to set up certificates. Please refer to the **Certificate
management** section.

Capture some traffic
~~~~~~~~~~~~~~~~~~~~

First put Hoverfly into *capture mode*. There are four ways to do this.

**A:** Use hoverctl to set the mode of the running Hoverfly:

::

    hoverctl mode capture

**B:** Start Hoverfly with the ``-capture`` flag:

::

    ./hoverfly -capture

**C:** Ensure Hoverfly is running (in any mode), then select "capture"
in the Admin UI, which is available at ``http://${HOVERFLY_HOST}:8888``.

**D:** Ensure Hoverfly is running (in any mode), then make an API call:

::

    curl -H "Content-Type application/json" -X PUT -d '{"mode":"capture"}' http://${HOVERFLY_HOST}:8888/api/v2/hoverfly/mode

The simplest way to capture some traffic is to use cURL:

::

    curl http://<MY_EXTERNAL_SERVICE> --proxy http://${HOVERFLY_HOST}:8500/

Or if you are developing or testing an application that uses an external
service, just run your tests.

To verify that Hoverfly is capturing traffic, you can look at the
Hoverfly logging output, or the Admin UI.

You can also make an API call to view captured traffic:

::

    curl http://${HOVERFLY_HOST}:8888/api/v2/simulation

What is happening?
~~~~~~~~~~~~~~~~~~

Hoverfly is transparently passing requests from the client application
through to the destination service, then passing the responses back. It
is storing the request/response pairs in memory, and persisting them on
disk in a file named ``requests.db``.

Deleting all captured traffic
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This can be done in three ways.

**A:** Use hoverctl to wipe Hoverfly

::

    hoverctl wipe

**B:** Visit the Admin UI at ``http://${HOVERFLY_HOST}:8888`` and click
"Wipe records"

**C:** Make an API call:

::

    curl -X DELETE http://${HOVERFLY_HOST}:8888/api/v2/simulation

Next steps
~~~~~~~~~~

Now you can use Hoverfly to mimic the external service. Proceed to the
**Simulating services** section.

If you want to export the traffic you have captured, proceed to the
**Exporting and importing** section.

Further reading
~~~~~~~~~~~~~~~

A detailed step-by-step guide to capturing traffic and creating a
simulated service is available here:

`Speeding up your slow
dependencies <https://specto.io/blog/speeding-up-your-slow-dependencies.html>`__
\*\*\* ## Simulating services ### Set up

To simulate a service, you will need to have either captured some
traffic (see the **Capturing traffic** section) or imported a service
data file (see the **Exporting and importing** section).

Simulate a service
~~~~~~~~~~~~~~~~~~

With Hoverfly running as a proxy
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

You will need to have your application or OS configured to use Hoverfly
as a proxy (see the **Installation and setup** section).

Put Hoverfly into *simulate mode*. There are four ways of doing this:

**A:** Use hoverctl to start Hoverfly and set the mode to simulate:

::

    hoverctl start
    hoverctl mode simulate

**B:** Start Hoverfly without specifying the mode (Hoverfly starts in
*simulate mode* by default)

::

    ./hoverfly

**C:** Ensure Hoverfly is running (in any mode), then select "simulate"
in the Admin UI, which is available at ``http://${HOVERFLY_HOST}:8888``
by default.

**D:** Ensure Hoverfly is running (in any mode), then make an API call:

::

    curl -H "Content-Type application/json" -X POST -d '{"mode":"simulate"}' http://${HOVERFLY_HOST}:8888/api/v2/hoverfly/mode

Provided you have previously captured traffic either using cURL or by
running tests, you can now repeat the steps you took to capture traffic.

Instead of using the external service you captured, your application (or
cURL) will now be using the Hoverfly simulation.

If you are running a test suite, you will probably notice that it runs
much faster.

With Hoverfly running as a webserver
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Put Hoverfly into *simulate mode* while it is running as a webserver.
There are two ways to do this:

**A:** Use hoverctl to start Hoverfly as a webserver, and set the mode
to simulate:

::

    hoverctl start webserver
    hoverctl mode simulate

**B:** Start Hoverfly with the ``-webserver`` flag and do not specify a
mode (Hoverfly starts in *simulate mode* by default)

::

    ./hoverfly -webserver

**NOTE:** When running as a webserver, Hoverfly currently only supports
*simulate* mode.

**NOTE:** Currently HTTPS is not supported when running Hoverfly as a
webserver. HTTPS support when running as a webserver is on the roadmap.

Provided you have previously captured traffic either using cURL or by
running tests with Hoverfly running as a proxy (or you have imported
previously captured traffic - see the **Exporting and importing**
section), Hoverfly can now be used instead of the external service you
captured.

Since Hoverfly is simulating the service as a *webserver* rather than a
*proxy*, you will need to substitute the host name of the service you
captured with the host and port that the Hoverfly webserver is running
on when you make a request.

For example, if you created a simulation by running the following
command with Hoverfly running in *capture* mode as a **proxy**:

::

    curl http://my.host.com/api/health

You would need to run the following command with Hoverfly running in
*simulate* mode as a **webserver** to get the captured response:

::

    curl http://${HOVERFLY_HOST}:8500/api/health

What is happening?
~~~~~~~~~~~~~~~~~~

Hoverfly is no longer passing requests through to the external service.
Instead, for each request it receives, it is returning a matched
response from its cache. The cache is stored in memory, and persisted on
disk in the ``requets.db`` file.

Next steps
~~~~~~~~~~

Now you have simulated a service, you can use middleware to simulate
network latency or failure, or manipulate data in the responses. Proceed
to the **Middleware** section.

Alternatively, you can learn how to use *synthesize mode* to generate
responses to requests on the fly. Proceed to the **Creating synthetic
services** section.

Further reading
~~~~~~~~~~~~~~~

A detailed step-by-step guide to capturing traffic and creating a
simulated service is available here:

`Speeding up your slow
dependencies <https://specto.io/blog/speeding-up-your-slow-dependencies.html>`__

A guide on how to capture and simulate the Meetup API is available here:

`Virtualizing the Meetup
API <https://specto.io/blog/hoverfly-meetup-api.html>`__

--------------

Simulating service latency
--------------------------

Once you have created a simulated service by capturing traffic between
your application and an external service, you may wish to make the
simulation more "realistic" by applying latency to the responses
returned by Hoverfly.

This could be done using middleware (See the **Using middleware**
section). However, if you do not want to go to the effort of writing a
middleware script, you can use a JSON file to apply a set of fixed
response delays to the Hoverfly simulation.

This method is useful if Hoverfly is being used in a load test to
simulate an external service, and there is a requirement to simulate
external service latency. Under high load, the overhead of executing
middleware scripts will impact the performance of Hoverfly, making the
middleware approach to adding latency unsuitable.

Set up
~~~~~~

To simulate service latency, you will need to have created a simulation
by capturing traffic (see the **Capturing traffic** and **Simulating
services** sections).

Simulate latency
~~~~~~~~~~~~~~~~

JSON configuration
^^^^^^^^^^^^^^^^^^

To simulate latency, Hoverfly can be configured to apply a delay to
individual hosts or API endpoints in the simulation using a JSON
configuration file. This is done using a regular expression to match
against the URL, a delay value in milliseconds, and an optional HTTP
method value.

For example, to apply a delay of 2 seconds to all hosts in the
simulation:

::

    {
      "data": [
        {
          "urlPattern": ".",
          "delay": 2000
        }
      ]
    }

To apply a delay of 1 second to ``1.myhost.com`` and a delay of 2
seconds to ``2.myhost.com``:

::

    {
      "data": [
        {
          "urlPattern": "1\\.myhost\\.com",
          "delay": 1000
        },
        {
          "urlPattern": "2\\.myhost\\.com",
          "delay": 2000
        }
      ]
    }

It is also possible to apply delays to specific resources and endpoints
in your API. In the following example, a delay of 1 second is applied to
all endpoints of resource A. For resource B, a delay of 1 second is
applied to the GET endpoint, but a different delay of 2 seconds is
applied to the POST endpoint:

::

    {
      "data": [
        {
          "urlPattern": "myhost\\.com\\/A",
          "delay": 1000
        },
        {
          "urlPattern": "myhost\\.com\\/B",
          "delay": 1000,
          "httpMethod": "GET"
        },
        {
          "urlPattern": "myhost\\.com\\/B",
          "delay": 2000,
          "httpMethod": "POST"
        }
      ]
    }

The **delays will be matched in the order that they appear in the JSON
configuration file**. In the following example, ``"urlPattern":"."``
matches all hosts, overriding ``"urlPattern": "1\\.myhost\\.com"`` and
all subsequent matches, applying a 3 second delay to all responses:

::

    {
      "data": [
        {
          "urlPattern": ".",
          "delay": 3000
        },
        {
          "urlPattern": "1\\.myhost\\.com",
          "delay": 1000
        }
      ]
    }

Applying the configuration
^^^^^^^^^^^^^^^^^^^^^^^^^^

The configuration can be applied using hoverctl.

To apply delays:

::

    hoverctl delays path/to/my_delays.json

To view the delays which have been applied:

::

    hoverctl delays

Alternatively, the configuration can be applied using the Hoverfly API
directly:

To apply delays:

::

    curl -H "Content-Type application/json" -X PUT -d '{"data":[{"urlPattern":"1\\.myhost\\.com","delay":1000},{"urlPattern":"2\\.myhost\\.com","delay":2000}]}' http://${HOVERFLY_HOST}:8888/api/delays

To view the delays which have been applied

::

    curl http://${HOVERFLY_HOST}:8888/api/delays

--------------

Managing simulation data
------------------------

Hoverfly can export and import service data in JSON format. This is
useful if:

-  You have captured some traffic and want to store it somewhere other
   than the Hoverfly ``requests.db`` file - in a Git repository for
   example.
-  You are running Hoverfly in a Docker container - so persisting data
   on disk is not ideal
-  You want to capture traffic, then modify it somehow before
   re-importing it
-  You want to share your service data someone else

Exporting captured data
~~~~~~~~~~~~~~~~~~~~~~~

Using hoverctl, you can export all the simulation data from Hoverfly.
Exported data will be written to a JSON file in your current working
directory.

::

    hoverctl export mysimulation.json

For more information about hoverctl, check the **Reference** section.

Simulation data JSON format
~~~~~~~~~~~~~~~~~~~~~~~~~~~

Hoverfly stores captured **Request Response Pairs** (i.e. "traffic") in
the following JSON structure:

::

    {
        "data": {
            "pairs": [{
                "response": {
                    "status": 200,
                    "body": "body here",
                    "encodedBody": false,
                    "headers": {
                        "Content-Type": ["text/html; charset=utf-8"]
                    }
                },
                "request": {
                    "requestType": "recording",
                    "path": "/",
                    "method": "GET",
                    "destination": "myhost.io",
                    "scheme": "https",
                    "query": "",
                    "body": "",
                    "headers": {
                        "Accept": ["text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"],
                        "User-Agent": ["Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"]
                    }
                }
            }]
        }
    }

When you export a simulation that you have captured, the JSON file will
look something like this. Notice that by default, the request
``requestType`` is ``recording``.

Base64 encoding binary data
^^^^^^^^^^^^^^^^^^^^^^^^^^^

As JSON does not support binary data, binary responses are base64
encoded. This is denoted by the ``encodedBody`` field. Hoverfly
automatically encodes and decodes the data during the export and import
processes.

A note about matching
^^^^^^^^^^^^^^^^^^^^^

Hoverfly works by inspecting requests being made, extracting key pieces
of information and then matching them against stored requests. Standard
matching uses everything in the request apart from headers. (This
because request headers often change depending on browser, HTTP client
and or the time of day.)

In some cases, you may want to use partial matching. For example, you
may want Hoverfly to return a specific response for **any** incoming
request going to a specific path. This can be achieved using request
templates.

Request templates (for partial matching)
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

If no exact match is found for an incoming request, Hoverfly will
attempt to match on request templates.

Request templates are defined in the JSON file by setting the
``requestType`` property for a request to ``template`` and including
**only** the information in the request that you want Hoverfly to use in
the match.

For example:

::

    {
        "data": {
            "pairs": [{
                "response": {
                    "status": 200,
                    "body": "<h1>Matched on template</h1>",
                    "encodedBody": false,
                    "headers": {
                        "Content-Type": ["text/html; charset=utf-8"]
                    }
                },
                "request": {
                    "requestType": "template",
                    "path": "/template"
                }
            }]
        }
    }

Here, any request with the path ``/template`` will return the same
response.

For looser matching, it is possible to use a wildcard to substitute
characters. This is achieved by using an ``*`` symbol. This will match
any number of characters, and is case sensitive.

It is possible to combine the wildcard (``*``) with characters to
substitute parts of a string. In the next example, we use a wildcard the
replace part of a URL path. This allows us to match on either
``/api/v1/template`` or ``/api/v2/template``.

::

    {
        "data": {
            "pairs": [{
                "response": {
                    "status": 200,
                    "body": "<h1>Matched on template</h1>",
                    "encodedBody": false,
                    "headers": {
                        "Content-Type": ["text/html; charset=utf-8"]
                    }
                },
                "request": {
                    "requestType": "template",
                    "path": "/api/*/template"
                }
            }]
        }
    }

The JSON file can contain both requests recordings and request
templates:

::

    {
        "data": {
            "pairs": [{
                "response": {
                    "status": 200,
                    "body": "<h1>Matched on recording</h1>",
                    "encodedBody": false,
                    "headers": {
                        "Content-Type": [
                            "text/html; charset=utf-8"
                        ]
                    }
                },
                "request": {
                    "requestType": "recording",
                    "path": "/",
                    "method": "GET",
                    "destination": "myhost.io",
                    "scheme": "https",
                    "query": "",
                    "body": "",
                    "headers": {
                        "Accept": [
                            "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
                        ],
                        "Content-Type": [
                            "text/plain; charset=utf-8"
                        ],
                        "User-Agent": [
                            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"
                        ]
                    }
                }
            }, {
                "response": {
                    "status": 200,
                    "body": "<h1>Matched on template</h1>",
                    "encodedBody": false,
                    "headers": {
                        "Content-Type": [
                            "text/html; charset=utf-8"
                        ]
                    }
                },
                "request": {
                    "requestType": "template",
                    "path": "/template",
                    "method": null,
                    "destination": null,
                    "scheme": null,
                    "query": null,
                    "body": null,
                    "headers": null
                }
            }],
            "globalActions": {
                "delays": []
            }
        }
    }

A standard workflow might be:

1. Capture some traffic
2. Export it to JSON
3. Edit the JSON to set certain requests to templates, removing the
   properties for these requests that should be excluded from the match
4. Re-import the JSON to Hoverfly

If the ``requestType`` property is not defined or not recognized,
Hoverfly will treat a request as a recording.

--------------

Importing simulation data
-------------------------

Service data can be imported from the local file system, or from a URL.
After you have edited captured data to include some request templates,
you will probably want to import it back into Hoverfly. There are four
ways to do this.

**A:** Use hoverctl to load a file into Hoverfly:

::

    hoverctl import simulation.json

For more about hoverctl, `check here <../reference/hoverctl.md>`__.

**B:** Start Hoverfly in *simulate mode* with the ``-import`` flag:

::

    ./hoverfly -import path/to/data.json
    ./hoverfly -import https://<MY_HOST>/data.json

Multiple service data files can be imported like this:

::

    ./hoverfly -import path/to/data_1.json -import path/to/data_2.json

If the file you specified cannot be found, Hoverfly will not start.

**C:** Use the ``HoverflyImport`` environment variable:

::

    export HoverflyImport="path/to/data.json"
    export HoverflyImport="https://<MY_HOST>/data.json"

If the file you specified cannot be found, Hoverfly will not start.

**D:** Make an API call:

::

    curl --data "@/path/to/data.json" http://${HOVERFLY_HOST}:8888/api/v2/simulation

--------------

Using middleware
----------------

Hoverfly middleware can be written in any language. Middleware modules
receive a service data JSON string via the standard input (STDIN) and
must return a service data JSON string to the standard output (STDOUT).

The way middleware is applied to the requests and/or the responses
depends on the mode:

-  Capture: Calls middleware once with an outgoing request object
-  Simulate: Calls middleware once with both request and response
   objects but any modifications will only affect the response
-  Synthesize: Calls middleware once with a request object so that the
   middleware can create a response
-  Modify Mode: Calls middleware twice, the first with a request object
   so that the request can be modified and then a second time with both
   the request and the response

**In *simulate mode*, middleware is only executed when a request is
matched to a stored response**.

To implement more dynamic behaviour, middleware should be combined with
Hoverfly's *synthesize mode* (see the **Creating synthetic services**
section).

Simple middleware examples
~~~~~~~~~~~~~~~~~~~~~~~~~~

Please read the **Installation and setup**, **Capturing traffic** and
**Simulating services** sections (if you haven't already) before
proceeding.

Middleware is set using hoverctl.

::

    hoverctl middleware "path/to/script"

You can also set the middleware on start up with the Hoverfly binary

::

    ./hoverfly -middleware "path/to/script"

This will start Hoverfly in the default mode (*simulate mode*) and
import a middleware script.

The string supplied as middleware can contain commands as well as a path
to a file. For example, if you have written middleware in Go:

::

    ./hoverfly -middleware "go run path/to/file.go"

::

    hoverctl middleware "go run path/to/file.go"

Python example
^^^^^^^^^^^^^^

This example will change the response code and body in each response,
and add 2 seconds of delay (simulating network latency).

Ensure that you have captured some traffic with Hoverfly.

Ensure that you have Python installed.

Save the following code into a file named ``example.py`` and make it
executable (``chmod +x example.py``):

::

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

Retart Hoverfly in *simulate mode* with the ``example.py`` script
specified as middleware:

::

    hoverctl middleware "./example.py"

Repeat the steps you took to capture the traffic.

You will notice that every response will have the ``201`` status code,
and the body will have been replaced by the string specified in the
script.

There will also be a 2 second delay between each request and the
response.

Javascript example
^^^^^^^^^^^^^^^^^^

This example will change the response code and body in each response.

Ensure that you have captured some traffic with Hoverfly.

Ensure that you have NodeJS installed.

Save the following code into a file named ``example.js`` and make it
executable (``chmod +x example.js``):

::

    #!/usr/bin/env node

    process.stdin.resume();  
    process.stdin.setEncoding('utf8');  
    process.stdin.on('data', function(data) {
      var parsed_json = JSON.parse(data);
      // changing response
      parsed_json.response.status = 201;
      parsed_json.response.body = "body was replaced by JavaScript middleware\n";

      // stringifying JSON response
      var newJsonString = JSON.stringify(parsed_json);

      process.stdout.write(newJsonString);
    });

Restart Hoverfly in *simulate mode* with the ``example.js`` script
specified as middleware:

::

    hoverctl middleware "./example.js"

Repeat the steps you took to capture the traffic.

You will notice that every response will have the ``201`` status code,
and the body will have been replaced by the string specified in the
script.

Remote execution of middleware
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The examples above show how to execute middleware on the machine on
which Hoverfly is running. It is also possible to execute middleware
remotely over HTTP.

For example, you could write a piece of middleware using `AWS
Lambda <https://docs.aws.amazon.com/lambda/latest/dg/welcome.html>`__
and then provide the URL to Hoverfly as middleware.

::

    hoverctl middleware "https://1sfefc4.execute-api.eu-west-1.amazonaws.com/remote/middleware"

Remote middleare execution works in a similar way to local middleware
execution. Hoverfly will send the same JSON data to the middleware, but
instead of sending it via ``stdin``, it is sent as a POST request to the
specified URL. The middleware is then expected to return the same data
as it would to ``stdout`` to the response body.

Base64 encoding binary data
~~~~~~~~~~~~~~~~~~~~~~~~~~~

The above example assumes that the responses do not contain binary data.
If they do then the ``encodedBody`` field will be set to ``true``. In
this circumstance, if you want to mutate the body you must base64 encode
and decode it.

It is also worth bearing in mind that the "Content-Length" header must
be set to the length of the unencoded body.

Further reading
~~~~~~~~~~~~~~~

A detailed step-by-step guide to using middleware is available here:

`Modifying traffic on the
fly <https://specto.io/blog/service-virtualization-is-so-last-year.html>`__

More middleware examples are available here:

`Middleware
examples <https://github.com/SpectoLabs/hoverfly/tree/master/examples/middleware>`__

Currently, middleware is not supported by hoverctl.

--------------

Creating synthetic services
---------------------------

In *synthesize mode*, Hoverfly does not make use of its cache.
*Synthesize mode* is dependent on middleware. For each request that
Hoverfly receives, it executes middleware, which must generate a
response. Hoverfly then returns the generated response to the client.

This mode allows you use Hoverfly as a dynamic stub server. For example,
you could use it with a set of template responses stored on a file
system or in a database, which could be populated with data and returned
based on any characteristic in the request.

TODO

--------------

Modifying traffic
-----------------

In *modify mode*, Hoverfly does not make use of its cache. Traffic is
passed through Hoverfly, and middleware is executed both on out-going
and in-bound traffic. In this mode, Middleware can be used to manipulate
any part of a request or a response.

This is a mode with potential applications outside of simulating
external services for development and testing. Possible use-cases
include introducing transparent redirects and injecting authentication
headers.

Set up
~~~~~~

You will need to have your application or OS configured to use Hoverfly
as a proxy (see the **Installation and setup** section).

To modify traffic on the fly, you do not need to have captured any
traffic or imported any Hoverfly JSON. However, it is important that you
read the **Middleware** section first, if you haven't already.

For the example below, you will also need Python installed on the
Hoverfly host.

Modify some traffic
~~~~~~~~~~~~~~~~~~~

Create a middleware script
^^^^^^^^^^^^^^^^^^^^^^^^^^

First, create a middleware script to modify traffic. Save the following
code in a file named ``modify.py`` in your Hoverfly directory and make
it executable with ``chmod +x modify.py``:

::

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

This middleware script will transparently redirect **any** request to
**any host** that is passed through Hoverfly to
``http://mirage.readthedocs.org``.

Put Hoverfly into Modify mode
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

There are three ways of doing this:

**A:** Start Hoverfly with the ``-modify`` flag and specify middleware
with the ``-middleware`` flag:

::

    ./hoverfly -modify -middleware "modify.py"

**B:** Ensure Hoverfly is running (in any mode), and that the middleware
has been specified, then select "modify" in the Admin UI which is
available at ``http://${HOVERFLY_HOST}:8888``.

**C:** Ensure Hoverfly is running (in any mode), and that the middleware
has been specified, then make an API call:

::

    curl -H "Content-Type application/json" -X POST -d '{"mode":"modify"}' http://${HOVERFLY_HOST}:8888/api/v2/hoverfly/mode

Now cURL any (HTTP) URL with Hoverfly as a proxy. For example:

::

    curl http://flask.readthedocs.org/ --proxy http://${HOVERFLY_HOST}:8500/

You will see that the response comes from
``http://mirage.readthedocs.org``.

What is happening?
~~~~~~~~~~~~~~~~~~

Hoverfly is executing the ``modify.py`` script on the request and the
response. The middleware takes the service data JSON string via STDIN,
and replaces the ``destination`` with the string
``mirage.readthedocs.org``. It also takes the response JSON string via
STDIN and replaces the ``method`` with ``201``.

Further reading
~~~~~~~~~~~~~~~

This example is taken from a more detailed step-by-step guide:

`Modifying traffic on the
fly <https://specto.io/blog/service-virtualization-is-so-last-year.html>`__

--------------

Filtering destination
---------------------

You may wish to control what Hoverfly captures or simulates. By default,
Hoverfly will process everything.

To specify which URLs Hoverfly processes you can use the
``hoverctl destination`` command. You can provide this either an exact
match or regex. This destination will be compared against the host and
the path of a URL. For example, we can tell Hoverfly to ignore anything
that isn't hoverfly.io.

::

    hoverctl destination "hoverfly.io"

But we could specficially say we are only interested in capturing any
request response if it contains an API.

::

    hoverctl destination "api"

Thhis destination would match on ``api.hoverfly.io/endpoint`` and
``hoverfly.io/api/endpoint``.

Hoverfly will only process those requests which match the specified
destination. All other requests will be passed through. This applies to
all modes. With a destination set, it is possible to request real
responses alongside simulated responses with simulate mode. \*\*\*

HTTPS support & certificate management
--------------------------------------

Hoverfly ships with a default certificate (``cert.pem`` in the
repository root directory). To use Hoverfly with HTTPS traffic, you will
need to add this default certificate to your trust store.

Using Hoverfly, you may also generate new certificates.

Generating and using certificates
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

To enable support for HTTPS services, Hoverfly can generate public and
private keys. To generate a key pair, use the ``-generate-ca-cert``
flag:

::

    ./hoverfly -generate-ca-cert

This will create ``cert.pem`` and a ``key.pem`` files in your current
directory. Next time you run Hoverfly, you can tell it to use these
certificate and key files using the ``-cert`` and ``-key`` flags:

::

    ./hoverfly -cert cert.pem -key key.pem

You can also define the certificate and key to use with Hoverctl.

::

    hoverctl start --certificate cert.pem --key.pem

Once Hoverfly has started with the new certificate and key file, you
will then need to add the ``cert.pem`` file to your trusted
certificates. Alternatively, you can turn off certificate verification.
For example, to make insecure requests with cURL, you can use the ``-k``
flag:

::

    curl https://www.bbc.co.uk --proxy http://${HOVERFLY_HOST}:8500 -k

Turn off verification when capturing or modifying traffic
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You can tell Hoverfly to ignore untrusted certificates when capturing or
modifying traffic in two ways.

**A:** Use the ``-tls-verification=false`` flag on startup:

::

    ./hoverfly -tls-verification=false

**B:** Use the --disable-tls flag on hoverctl

::

    hoverctl start --disable-tls

--------------

Authentication
--------------

Hoverfly uses a combination of Basic Auth and `JWT <https://jwt.io/>`__
(JSON Web Tokens) to authenticate users. Authentication is disabled by
default.

Enabling authentication
~~~~~~~~~~~~~~~~~~~~~~~

If you enable authentication, and you haven't created a user using flags
or environment variables (see below), you will be prompted to create a
new user when you start Hoverfly.

To enable authentication, you can use the ``-auth`` flag on startup:

::

    ./hoverfly -auth

Or you can use the ``HoverflyAuthDisabled`` environment variable:

::

    export HoverflyAuthEnabled=true

If the ``-auth`` flag is supplied **or** the ``HoverflyAuthEnabled``
environment variable is set to ``true``, authentication will be enabled.

When authentication is disabled, **any username and password
combination** can be used to access the Admin UI.

Adding users
~~~~~~~~~~~~

You can add a user using the ``-add``, ``-username`` and ``-password``
flags at startup:

::

    ./hoverfly -add -username <username> -password <password>

This will add an admin user. To add a non-admin user, use the ``-admin``
flag:

::

    ./hoverfly -add -username <username> -password <password> -admin false

You can also add an initial super user using environment variables. This
is useful if you are using Hoverfly in Docker, for example:

::

    export HoverflyAdmin="username"
    export HoverflyAdminPass="password"

Token usage for API authentication
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

To get the token for a user, make an API call:

::

    curl -H "Content-Type application/json" -X POST -d '{"Username": "<username>", "Password": "<password>"}' http://${HOVERFLY_HOST}:8888/api/token-auth

To use the token in an API call:

::

    curl -H "Authorization: Bearer <token>" http://${HOVERFLY_HOST}:8888/api/v2/simulation

By default, tokens expire after one day. You can override this by
setting the ``HoverflyTokenExpiration`` environment variable in seconds:

::

    export HoverflyTokenExpiration=3600

Setting the Hoverfly secret
~~~~~~~~~~~~~~~~~~~~~~~~~~~

By default, a new random secret will be generated every time you launch
Hoverfly. However, you can specify a secret using the ``HoverflySecret``
environment variable:

::

    export HoverflySecret=<my_secret>

