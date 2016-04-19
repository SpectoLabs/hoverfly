# Capturing traffic

## Set up

To capture traffic with Hoverfly, you will need Hoverfly on your local machine and your OS or application will need to be configured to use Hoverfly as a proxy. 

Follow the [Getting Started](#) guide to do this.

If the external service you wish to capture requires HTTPS, [you will need to configure certificates](#). 

## Set capture mode

Hoverfly will now be able to intercept requests between your application and the external service you want to capture.

There are three ways to put Hoverfly into capture mode.
 
You can start Hoverfly using the `-capture` flag:
 
    ./hoverfly -capture

You can select "Capture" in the [AdminUI](#) at `http://<HOVERFLY_HOST>:8888`

Or you can use an [API](#) call: 

    curl -H "Content-Type application/json" -X POST -d '{"mode":"capture"}' http://<HOVERFLY_HOST>:8888/api/state


## Capture some traffic

If you have some tests that use an external service, and you have your environment configured to use Hoverfly as a proxy, you can now run your tests. Hoverfly will capture all the requests your application makes and the responses returned by the external service.

Alternatively you can just use cURL to capture some traffic:

    curl http://<my_external_service> --proxy http://localhost:8500/
    
Look at the Hoverfly AdminUI to see the traffic being captured.
    
Now you have captured some traffic, you can [switch Hoverfly into "virtualize" mode](#) and run your tests again - or execute the same cURL command. Your application (or cURL) will now be talking to Hoverfly. 
   
A [more detailed step-by-step guide to capturing and virtualizing traffic is available here](https://specto.io/blog/speeding-up-your-slow-dependencies.html).    