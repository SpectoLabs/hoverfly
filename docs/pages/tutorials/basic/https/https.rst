.. _simulating_https:

Simulating HTTPS APIs
=====================

To capture HTTPS traffic, you need to use Hoverfly's SSL certificate.

First, download the certificate:

::

    wget https://raw.githubusercontent.com/SpectoLabs/hoverfly/master/core/cert.pem

We can now run Hoverfly with the standard ``capture`` then ``simulate`` workflow.

::

    hoverctl start
    hoverctl mode capture
    curl --proxy localhost:8500 https://example.com --cacert cert.pem
    hoverctl mode simulate
    curl --proxy localhost:8500 https://example.com --cacert cert.pem
    hoverctl stop


Curl makes the HTTPS request by first establishing a TLS tunnel to the destination with Hoverfly using the HTTP CONNECT method.
As curl has supplied Hoverfly's SSL certificate, Hoverfly is then able to intercept and capture the traffic.
Effectively SSL-encrypted communication (HTTPS) is established through an unencrypted HTTP proxy.

.. note::

  This example uses cURL. If you are using Hoverfly in another environment, you will need to add the certificate to your trust store.
  This is done automatically by the Hoverfly Java library (see :ref:`hoverfly_java`).

.. seealso::
  
   This example uses Hoverfly's default SSL certificate. Alternatively, you can use Hoverfly to generate
   a new certificate. For more information, see :ref:`configuressl`. 
