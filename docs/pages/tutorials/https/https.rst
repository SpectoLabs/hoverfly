Simulating HTTPS APIs
---------------------

Simulating HTTPS traffic is rather similar to simulating HTTP traffic, with the addition of having to properly add / configure our HTTPS certificate.

Let's begin by downloading Hoverfly's certificate:

::

    wget https://raw.githubusercontent.com/SpectoLabs/hoverfly/master/core/cert.pem

We can now run hoverfly, with a standard ``capture`` then ``simulate`` workflow.

::

    hoverctl start
    hoverctl mode capture
    curl --proxy https://localhost:8500 https://example.com --cacert cert.pem
    hoverctl mode simulate
    curl --proxy https://localhost:8500 https://example.com --cacert cert.pem
    hoverctl stop

Curl was able to make the https request, using an https proxy, since we provided it with Hoverfly's ssl certificate.