.. figure:: logo-large.png
   :alt: Hoverfly logo

What is Hoverfly?
-----------------

Hoverfly is a lightweight, open source `service
virtualization <https://en.wikipedia.org/wiki/Service_virtualization>`__
tool. Using Hoverfly, you can virtualize your application dependencies
to create a self-contained development or test environment.

Hoverfly is a proxy written in Go. It can capture HTTP(s) traffic
between an application under test and external services, and then
replace the external services. It can also generate synthetic responses
on the fly.

--------------

Get Hoverfly
------------

Hoverfly is a single binary file. It comes with an optional command line
interface tool called hoverctl.

Download one of the zip files below, extract the Hoverfly and hoverctl
binaries, and move them to a directory on your
`PATH <https://www.java.com/en/download/help/path.xml>`__.

-  `macOS
   64bit <https://github.com/SpectoLabs/hoverfly/releases/download/v0.9.2/hoverfly_bundle_OSX_amd64.zip>`__
-  `Linux
   32bit <https://github.com/SpectoLabs/hoverfly/releases/download/v0.9.2/hoverfly_bundle_linux_386.zip>`__
-  `Linux
   64bit <https://github.com/SpectoLabs/hoverfly/releases/download/v0.9.2/hoverfly_bundle_linux_amd64.zip>`__
-  `Windows
   32bit <https://github.com/SpectoLabs/hoverfly/releases/download/v0.9.2/hoverfly_bundle_windows_386.zip>`__
-  `Windows
   64bit <https://github.com/SpectoLabs/hoverfly/releases/download/v0.9.2/hoverfly_bundle_windows_amd64.zip>`__

Use `Homebrew <http://brew.sh/>`__ to install Hoverfly and hoverctl:

::

    brew install SpectoLabs/tap/hoverfly

--------------

Run Hoverfly
------------

To capture traffic between your application and an external service, you
will need to configure your OS, browser or application to use Hoverfly
as a proxy.

MacOS and Linux
~~~~~~~~~~~~~~~

Run Hoverfly using hoverctl:

::

    hoverctl start

By default, the Hoverfly proxy runs on localhost:8500. Switch Hoverfly
to "capture" mode and make a request with cURL, using Hoverfly as a
proxy:

::

    hoverctl mode capture
    curl --proxy http://localhost:8500 http://hoverfly.io/

Hoverfly has captured the request and the response. View the Hoverfly
logs:

::

    hoverctl logs

Switch Hoverfly to "simulate" mode" and make the same request:

::

    hoverctl mode simulate
    curl --proxy http://localhost:8500 http://hoverfly.io/

Hoverfly has returned the captured response.

Windows
~~~~~~~

Open a command prompt and run Hoverfly using hoverctl:

::

    hoverctl start

Configure your application, browser or OS to use the Hoverfly proxy
(http://localhost:8500). Switch Hoverfly to "capture" mode:

::

    hoverctl mode capture

Make some requests from your application, browser or OS, then view the
Hoverfly logs:

::

    hoverctl logs

Switch Hoverfly to "simulate" mode:

::

    hoverctl mode simulate

Make the same requests from your browser, OS or application. Hoverfly is
returning the captured responses.

More information on proxy settings:

-  `Windows proxy settings explained <http://blog.raido.be/?p=426>`__
-  `Firefox proxy
   setting <https://support.mozilla.org/en-US/kb/advanced-panel-settings-in-firefox#w_connection>`__
-  `Java Networking and
   Proxies <https://docs.oracle.com/javase/6/docs/technotes/guides/net/proxies.html>`__

.. toctree::
   :maxdepth: 2
   :hidden:

   introduction/introduction
   usage/usage
   reference/reference

