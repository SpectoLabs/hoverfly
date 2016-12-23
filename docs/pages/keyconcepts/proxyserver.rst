Hoverfly as a Proxy Server
--------------------------

One of Hoverfly's modes is as a proxy server. As you may or may not know, a proxy server passes requests between a client and server.

.. figure:: proxyserver.mermaid.png

It is sometimes essential to go via a proxy server to reach a network, for example, as a security measure. Therefore all network-enabled software can be configured to use a proxy.

Proxy configurations don't have to be a one to one relationship, they can also be a one to many, many to one, or many to many relationship.

.. figure:: proxyconfigs.mermaid.png

By default Hoverfly starts as a proxy server.

Using a proxy server
~~~~~~~~~~~~~~~~~~~~

The way you would normally use a proxy server, is by setting environment variables:

.. code:: bash

    export HTTP_PROXY="http://proxy-address:port"
    export HTTP_PROXYS="https://proxy-address:port"

Launching network-enabled software within an environment containing these variables *should* make it point to that particular proxy server. The term *should* is used as not all software respects these environment variables for security reasons.

In cases where applications don't enforce reading these environment variables, they can usually be launched with extra flags. Curl is one of these applications.

.. code:: bash

    curl http://google.com --proxy http://proxy-ip:port

.. note::
    
    The methods hereby described for enabling a proxy server are only intended to help you get started with Hoverfly, i.e they're mostly to help you execute code examples within this documentation. Every operating system, and Application can have its own methods for setting, enabling, and using a proxy server.

The difference between a proxy server and a web server
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

A proxy server is a type of web server. The main difference is that when a web server recieves a request from a client, it is expected to respond to it with whatever the intended response is (an html page, for example). The data it responds with is also generally expected to reside on that server, or within the same network.

A proxy server however is expected to pass the request on to another server. It is also expected to set some appropriate headers along the way, such as `X-Forwarded-For <https://en.wikipedia.org/wiki/X-Forwarded-For>`_, `X-Real-IP <https://en.wikipedia.org/wiki/X-Real-IP>`_, `X-Forwarded-Proto <https://en.wikipedia.org/wiki/X-Forwarded-Proto>`_ etc. Once a response is recieved from the destination, the proxy server is expected to pass it back to the client.

.. raw:: html
    
    <style>
        img {
            max-width: 600px !important;
        }
    </style>
