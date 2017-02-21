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
    curl --proxy https://localhost:8500 https://example.com --cacert cert.pem
    hoverctl mode simulate
    curl --proxy https://localhost:8500 https://example.com --cacert cert.pem
    hoverctl stop

Curl was able to make the HTTPS request using an HTTPS proxy because we provided it with Hoverfly's SSL certificate.

.. note::

  This example uses cURL. If you are using Hoverfly in another environment, you will need to add the certificate to your trust store.
  This is done automatically by the Hoverfly Java library (see :ref:`hoverfly_java`).

.. seealso::
  
   This example uses Hoverfly's default SSL certificate. Alternatively, you can use Hoverfly to generate
   a new certificate. For more information, see :ref:`configuressl`. 