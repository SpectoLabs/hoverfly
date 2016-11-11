# Capturing traffic

## Set up

To capture traffic, you will need:

* A client application (this could be an application you are developing or testing, or cURL if you just want to experiment)
* An external service or API that your client application communicates with over HTTP or HTTPS
* Hoverfly - either as a binary or the Docker image
* Your application or OS configured to use Hoverfly as an HTTP or HTTPS proxy

Please follow the steps in the **Installation and setup** section if you haven't already.

If you want to capture traffic to and from a service over HTTPS, you will need to set up certificates. Please refer to the **Certificate management** section.

## Capture some traffic

First put Hoverfly into *capture mode*. There are four ways to do this.

1. Use hoverctl to set the mode of the running Hoverfly:

        hoverctl mode capture

2. Start Hoverfly with the `-capture` flag:

        ./hoverfly -capture

3. Ensure Hoverfly is running (in any mode), then select "capture" in the Admin UI, which is available at `http://${HOVERFLY_HOST}:8888`.

4. Ensure Hoverfly is running (in any mode), then make an API call:

        curl -H "Content-Type application/json" -X PUT -d '{"mode":"capture"}' http://${HOVERFLY_HOST}:8888/api/v2/hoverfly/mode

The simplest way to capture some traffic is to use cURL:

    curl http://<MY_EXTERNAL_SERVICE> --proxy http://${HOVERFLY_HOST}:8500/
    
Or if you are developing or testing an application that uses an external service, just run your tests.

To verify that Hoverfly is capturing traffic, you can look at the Hoverfly logging output, or the Admin UI.

You can also make an API call to view captured traffic:

    curl http://${HOVERFLY_HOST}:8888/api/v2/simulation

## What is happening?

Hoverfly is transparently passing requests from the client application through to the destination service, then passing the responses back. It is storing the request/response pairs in memory, and persisting them on disk in a file named `requests.db`. 

## Deleting all captured traffic

This can be done in three ways.

1. Use hoverctl to wipe Hoverfly

       hoverctl wipe
       
2. Visit the Admin UI at `http://${HOVERFLY_HOST}:8888` and click "Wipe records"

3. Make an API call:

       curl -X DELETE http://${HOVERFLY_HOST}:8888/api/v2/simulation

## Next steps

Now you can use Hoverfly to mimic the external service. Proceed to the **Simulating services** section.

If you want to export the traffic you have captured, proceed to the **Exporting and importing** section.

## Further reading

A detailed step-by-step guide to capturing traffic and creating a simulated service is available here:

[Speeding up your slow dependencies](https://specto.io/blog/speeding-up-your-slow-dependencies.html)
