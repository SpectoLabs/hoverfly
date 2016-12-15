API
---

GET /api/v2/simulation
~~~~~~~~~~~~~~~~~~~~~~

Gets all simulation data from the running instance of Hoverfly. This
includes recordings, templates, delays and metadata.

Example response body
^^^^^^^^^^^^^^^^^^^^^

::

    {
      "data": {
        "pairs": [
          {
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
          },
          {
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
          }
        ],
        "globalActions": {
          "delays": []
        }
      },
      "meta": {
        "schemaVersion": "v1",
        "hoverflyVersion": "v0.9.0",
        "timeExported": "2016-11-11T11:53:52Z"
      }

--------------

PUT /api/v2/simulation
~~~~~~~~~~~~~~~~~~~~~~

Puts simulation into the running instance of Hoverfly, overwriting any
existing simulation data.

Example request body
^^^^^^^^^^^^^^^^^^^^

::

    {
      "data": {
        "pairs": [
          {
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
          },
          {
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
          }
        ],
        "globalActions": {
          "delays": []
        }
      },
      "meta": {
        "schemaVersion": "v1",
        "hoverflyVersion": "v0.9.0",
        "timeExported": "2016-11-11T11:53:52Z"
      }

--------------

GET /api/v2/hoverfly
~~~~~~~~~~~~~~~~~~~~

Gets configuration information from the running instance of Hoverfly.

Example response body
^^^^^^^^^^^^^^^^^^^^^

::

    {
        destination: ".",
        middleware: "",
        mode: "simulate",
        usage: {
            counters: {
                capture: 0,
                modify: 0,
                simulate: 0,
                synthesize: 0
            }
        }
    }

--------------

GET /api/v2/hoverfly/destination
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Gets the current destination setting for the running instance of
Hoverfly.

Example response body
^^^^^^^^^^^^^^^^^^^^^

::

    {
        destination: "."
    }

--------------

PUT /api/v2/hoverfly/destination
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Sets a new destination for the running instance of Hoverfly, overwriting
the existing destination setting.

Example request body
^^^^^^^^^^^^^^^^^^^^

::

    {
        destination: "new-destination"
    }

--------------

GET /api/v2/hoverfly/middleware
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Gets the middleware value for the running instance of Hoverfly. This
will be either an executable command, or an executable command with a
path to a middleware script.

Example response body
^^^^^^^^^^^^^^^^^^^^^

::

    {
        "middleware": "python ~/middleware.py"
    }

--------------

PUT /api/v2/hoverfly/middleware
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Sets a new middleware value, overwriting the existing middleware value
for the running instance of Hoverfly. The middleware value should be an
executable command, or an executable command with a path to a middleware
script. The command and the file must be available on the Hoverfly host
machine.

Example request body
^^^^^^^^^^^^^^^^^^^^

::

    {
        "middleware": "python ~/new-middleware.py"
    }

--------------

GET /api/v2/hoverfly/mode
~~~~~~~~~~~~~~~~~~~~~~~~~

Gets the mode for the running instance of Hoverfly.

Example response body
^^^^^^^^^^^^^^^^^^^^^

::

    {
        mode: "simulate"
    }

--------------

PUT /api/v2/hoverfly/mode
~~~~~~~~~~~~~~~~~~~~~~~~~

Changes the mode of the running instance of Hoverfly.

Example request body
^^^^^^^^^^^^^^^^^^^^

::

    {
        mode: "simulate"
    }

--------------

GET /api/v2/hoverfly/usage
~~~~~~~~~~~~~~~~~~~~~~~~~~

Gets metrics information for the running instance of Hoverfly.

Example response body
^^^^^^^^^^^^^^^^^^^^^

::

    {
        "metrics": {
            "counters": {
                "capture": 0,
                "modify": 0,
                "simulate": 0,
                "synthesize": 0
            }
        }
    }

--------------

Flags and environment variables
-------------------------------

Hoverfly can be configured using flags on startup, or using environment
variables.

Authentication
~~~~~~~~~~~~~~

Enable/disable authentication
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Flag:

::

    -auth <string>

Environment variable:

::

    export HoverflyAuthEnabled=<string>

Supply ``true`` to enable authentication. Defaults to ``false``.

Add a new user
^^^^^^^^^^^^^^

Flags:

::

    -add -username <string> -password <string> -admin <string>

Supply '-admin false' to make this a non-admin user (defaults to
'true').

For example:

::

    ./hoverfly -add -username username -password password -admin false      

This creates a new non-admin user with the username 'username' and the
password 'password'.

Environment variables:

::

    export HoverflyAdmin="username"
    export HoverflyAdminPass="password"

Setting these environment variables will create a new admin user when
Hoverfly starts.

Set Hoverfly secret
^^^^^^^^^^^^^^^^^^^

By default, a random secret is generated every time Hoverfly starts.

Environment variable:

::

    export HoverflySecret=<string>

Set API token expiration (in seconds)
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Set to one day by default.

Environment variable:

::

    export HoverflyTokenExpiration=<string>

Port selection
~~~~~~~~~~~~~~

Set the Admin UI/API port
^^^^^^^^^^^^^^^^^^^^^^^^^

Defaults to 8888.

Flag:

::

    -ap <string>

Environment variable:

::

    export AdminPort=<string>

Set the proxy port
^^^^^^^^^^^^^^^^^^

Defaults to 8500.

Flag:

::

    -pp <string>

Environment variable:

::

    export ProxyPort=<string>

Mode selection, import & middleware
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

By default, Hoverfly starts in *simulate mode*.

Set capture mode
^^^^^^^^^^^^^^^^

Flag:

::

    -capture

Set synthesize mode
^^^^^^^^^^^^^^^^^^^

Requires middleware to be specified.

Flag:

::

    -synthesize

Set modify mode
^^^^^^^^^^^^^^^

Requires middleware to be specified.

Flag:

::

    -modify

Specify middleware
^^^^^^^^^^^^^^^^^^

Flag:

::

    -middleware <string>

Supply the path to the middleware script.

For example:

::

    ./hoverfly -synthesize -middleware "scripts/gen_response.py"

Import service data
^^^^^^^^^^^^^^^^^^^

Flag:

::

    -import <string>

Import a service data JSON file from file system or URL. For example:

::

    ./hoverfly -import http://mypage.com/service_x.json

::

    ./hoverfly -import path/to/my/service_x.json      

Environment variable:

::

    export HoverflyImport=<string>

For example:

::

    export HoverflyImport="http://mypage.com/service_x.json"

Webserver
~~~~~~~~~

Turn Hoverfly into a simulation webserver
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Flag:

::

    -webserver

Destination
~~~~~~~~~~~

Specify which hosts to process
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Flag:

::

    -dest <string>

For example:

::

    ./hoverfly -dest fooservice.org -dest barservice.org -dest catservice.org

This will start Hoverfly in *simulate mode*, and only simulate requests
that are sent to fooservice.org, barservice.org and catservice.org.
Requests to all other hosts will pass through.

Specify host URI
^^^^^^^^^^^^^^^^

Use regular expression. Defaults to "."

Flag:

::

    -destination <string>

Persistence
~~~~~~~~~~~

Specify BoltDB or in-memory
^^^^^^^^^^^^^^^^^^^^^^^^^^^

Flag:

::

    -db <string>

By default, Hoverfly uses BoltDB to store data in a file on disk. Use
``-db memory`` to disable this and use in-memory persistence only.

Set BoltDB file
^^^^^^^^^^^^^^^

By default, a ``requests.db`` file is created in the Hoverfly directory.

Flag:

::

    -db-path <string>

Environment variable:

::

    export HoverflyDB=<string>

The file will be created if it doesn't exist.

TLS & Certificate management
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Generate certificate
^^^^^^^^^^^^^^^^^^^^

Hoverfly will generate private and public keys in the current directory.

Flags:

::

    -generate-ca-cert -cert-name <string> -cert-org <string>

Certificate name defaults to "hoverfly.proxy". Organization name
defaults to "Hoverfly Authority".

Use certificate and key
^^^^^^^^^^^^^^^^^^^^^^^

Supply paths to certificate and key file.

Flags:

::

    -cert <string> -key <string>

Turn off TLS verification
^^^^^^^^^^^^^^^^^^^^^^^^^

Defaults to "true".

Flag:

::

    -tls-verification=<string>

Environment variable:

::

    export HoverflyTlsVerification=<string>

Logging & metrics
~~~~~~~~~~~~~~~~~

Enable verbose mode
^^^^^^^^^^^^^^^^^^^

Logs every proxy request to STDOUT.

Flag:

::

    -v

Enable metrics logging
^^^^^^^^^^^^^^^^^^^^^^

Logs metrics to STDOUT.

Flag:

::

    -metrics

Misc
~~~~

Use uncompiled static files
^^^^^^^^^^^^^^^^^^^^^^^^^^^

Serve Admin UI static files directly from ./static/dist instead from
statik binary.

Flag:

::

    -dev

Get version
^^^^^^^^^^^

Get the version of Hoverfly

Flag:

::

    -version   

--------------

Hoverctl
--------

Hoverctl is a command line tool bundled with Hoverfly. The purpose of
hoverctl is to help in the managing of one or many instances of
Hoverfly. Hoverctl does not support all the functionality of Hoverfly
yet, but its feature set is growing.

.hoverfly directory
~~~~~~~~~~~~~~~~~~~

Hoverctl stores its state in a ``.hoverfly`` directory. Hoverctl will
create this directory in your home folder the first time it needs to
save state. This directory is used for the configuration for Hoverfly,
the process identifiers and the log files. Hoverctl will always check
the working directory before your home directory when looking for the
``.hoverfly`` directory. This allows for multiple configurations on a
per project basis if you require different configurations for Hoverfly.

.. code:: sh

    .hoverfly/config.yaml
    .hoverfly/hoverfly.8888.8500.pid
    .hoverfly/hoverfly.8888.8500.log

Configuration
^^^^^^^^^^^^^

::

    hoverfly.host             #default "localhost"
    hoverfly.admin.port       #default "8888"
    hoverfly.proxy.port       #default "8500"
    hoverfly.username         #default ""
    hoverfly.password         #default ""
    hoverfly.db.type          #default "memory"
    hoverfly.webserver        #default false
    hoverfly.tls.certificate" #default ""
    hoverfly.tls.key          #default ""
    hoverfly.tls.disable      #default false

Pid and log files
~~~~~~~~~~~~~~~~~

For each Hoverfly process created with hoverctl, one file is created to
store the process identifier and another for the STDOUT and STDERR of
Hoverfly. These files will be named using the hoverfly process with both
the admin and proxy ports.

Hoverctl commands
^^^^^^^^^^^^^^^^^

start
^^^^^

Hoverctl will let you start a Hoverfly process. For this to work, you
need to have the Hoverfly binary either in the same directory as
hoverctl or have Hoverfly on your $PATH. Hoverctl will start Hoverfly
based on the configuration defined in your ``config.yaml``. There is no
limit to the number of Hoverfly processes you can start. The only
requirement is that each Hoverfly process has its own unique admin and
proxy ports.

::

    hoverctl start

Using the global flags, it is possible to override the configuration
being set when starting an instance of Hoverfly.

::

    hoverctl start --disable-tls

By default, hoverctl will start Hoverfly as a proxy. If you wish to
start Hoverfly as a webserver instead:

::

    hoverctl start webserver

stop
^^^^

You can also stop Hoverfly processes using Hoverctl.

::

    hoverctl stop

The global flags can also be used here. If you have started an instance
of Hoverfly on a different admin and proxy port to your config.yaml, you
can still stop this by using the flags in combination with the stop
command.

::

    hoverctl stop --admin-port 1234 --proxy-port 4321

mode
^^^^

Using hoverctl, you can find out which mode Hoverfly is running in.

::

    hoverctl mode

You can also change the mode by specifying the name of mode you want
Hoverfly to be in.

::

    hoverctl mode capture

delete
^^^^^^

Hoverfly stores internal state while its running. This state is used for
testing your application. Using the delete command, you can specify what
you want to delete from Hoverfly.

::

    hoverctl delete simulations
    hoverctl delete delays
    hoverctl delete all

export
^^^^^^

Instead of having to save the response from the records API endpoint,
you can use the export function to save your simulation to disk.

::

    hoverctl export simulation.json

import
^^^^^^

Once you have simulations saved, you can import them into Hoverfly using
the import function.

::

    hoverctl import simulations.json

If your simulation file is hosted over HTTP, you can use hoverctl to
import it.

::

    hoverctl import http://example.org/simulation.json

If you have older, v1 simulations, you may still import them using the
``v1`` flag.

::

    hoverctl import --v1 old-simulations.json

delays
^^^^^^

If you want to apply delays to individual hosts in a simulation (to
simulate netwrok latency, for example), you can use the ``delays``
function to supply a JSON file containing the delay configuration or to
view delays which have been applied (See **Simulating service latency**
in the **Usage** section).

Set delays by supplying JSON file:

::

    hoverctl delays path/to/my_delays.json

Show delays that have been set:

::

    hoverctl delays

templates
^^^^^^^^^

As well importing request/response data using import, you can also
import request templates for partial matching to a response using the
``templates`` function. This function works with a JSON file containing
the JSON schema for request templates and responses. (See **Matching
requests** in the **Usage** section).

Set templates by supplying JSON file:

::

    hoverctl templates path/to/my_request_templates.json

Show templates that have been set:

::

    hoverctl templates

middleware
^^^^^^^^^^

This function is used for getting and setting the middleware being
executed by Hoverfly.

To get the middleware currently being used by Hoverfly

::

    hoverctl middleware

To set the middleware Hoverfly to use

::

    hoverctl middleware "middleware.sh"

The value given to the middleware function should be a string that
contains either a file path, a command a file path or a URL.

destination
^^^^^^^^^^^

This command is used for getting and setting the destination being used
to determine which requests are being processed by Hoverfly.

To get the destination currently being used by Hoverfly

::

    hoverctl destination

To set the destination value Hoverfly should use

::

    hoverctl destination 'hoverfly.io'

The value used should compile to Golang regex. Hoverctl will attempt to
compile the expression down and will warn the user if it does not
compile.

You can also test your destination value using the ``--dry-run`` flag.
This flag will not set the destination, but instead will test if your
regex pattern matches your intended response to record.

::

    hoverctl destination '\.io' --dry-run http://hoverfly.io

config
^^^^^^

This command is used for getting the file location of the config.yaml
being used. This command will also print the configuration that hoverctl
is using.

::

    hoverctl config

logs
^^^^

Used to get the logs from the instance of Hoverfly started with the
hoverctl. This command will return all the logs from when the process
was started

::

    hoverctl logs

If you are trying to debug what is happening and you need to watch the
Hoverfly logs, you can use the ``--follow`` flag to tail the logs and
watch them in real time.

::

    hoverctl logs --follow

Hoverctl flags
~~~~~~~~~~~~~~

--host
^^^^^^

This is a global flag that can be used to override the hoverfly.host
configuration value from the config.yaml file.

--admin-port
^^^^^^^^^^^^

This is a global flag that can be used to override the
hoverfly.admin.port configuration value from the config.yaml file.

--proxy-port
^^^^^^^^^^^^

This is a global flag that can be used to override the
hoverfly.proxy.port configuration value from the config.yaml file.

--certificate
^^^^^^^^^^^^^

This is a global flag that can be used to override the
hoverfly.tls.certificate configuration value from the config.yaml file.

--key
^^^^^

This is a global flag that can be used to override the hoverfly.tls.key
configuration value from the config.yaml file.

--disable-tls
^^^^^^^^^^^^^

This is a global flag that can be used to override the
hoverfly.tls.disable configuration value from the config.yaml file.

--verbose
^^^^^^^^^

This is a global flag that can be used to get the verbose logs from
hoverctl.

--version (-v)
^^^^^^^^^^^^^^

This is a global flag that can be used to get the version of hoverctl.
