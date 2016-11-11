# Simulating services
## Set up

To simulate a service, you will need to have either captured some traffic (see the **Capturing traffic** section) or imported a service data file (see the **Exporting and importing** section).

## Simulate a service
### With Hoverfly running as a proxy

You will need to have your application or OS configured to use Hoverfly as a proxy (see the **Installation and setup** section).

Put Hoverfly into *simulate mode*. There are four ways of doing this:

1. Use hoverctl to start Hoverfly and set the mode to simulate:

       hoverctl start
       hoverctl mode simulate

2. Start Hoverfly without specifying the mode (Hoverfly starts in *simulate mode* by default)

       ./hoverfly

3. Ensure Hoverfly is running (in any mode), then select "simulate" in the Admin UI, which is available at `http://${HOVERFLY_HOST}:8888` by default.

4. Ensure Hoverfly is running (in any mode), then make an API call:

       curl -H "Content-Type application/json" -X POST -d '{"mode":"simulate"}' http://${HOVERFLY_HOST}:8888/api/v2/hoverfly/mode

Provided you have previously captured traffic either using cURL or by running tests, you can now repeat the steps you took to capture traffic.

Instead of using the external service you captured, your application (or cURL) will now be using the Hoverfly simulation.

If you are running a test suite, you will probably notice that it runs much faster.

### With Hoverfly running as a webserver
Put Hoverfly into *simulate mode* while it is running as a webserver. There are two ways to do this:

1. Use hoverctl to start Hoverfly as a webserver, and set the mode to simulate:

    hoverctl start webserver
    hoverctl mode simulate

2. Start Hoverfly with the `-webserver` flag and do not specify a mode (Hoverfly starts in *simulate mode* by default)

       ./hoverfly -webserver

**NOTE:** When running as a webserver, Hoverfly currently only supports *simulate* mode.  

**NOTE:** Currently HTTPS is not supported when running Hoverfly as a webserver. HTTPS support when running as a webserver is on the roadmap.

Provided you have previously captured traffic either using cURL or by running tests with Hoverfly running as a proxy (or you have imported previously captured traffic - see the **Exporting and importing** section), Hoverfly can now be used instead of the external service you captured.

Since Hoverfly is simulating the service as a *webserver* rather than a *proxy*, you will need to substitute the host name of the service you captured with the host and port that the Hoverfly webserver is running on when you make a request.

For example, if you created a simulation by running the following command with Hoverfly running in *capture* mode as a **proxy**:

      curl http://my.host.com/api/health

You would need to run the following command with Hoverfly running in *simulate* mode as a **webserver** to get the captured response:

      curl http://${HOVERFLY_HOST}:8500/api/health


## What is happening?

Hoverfly is no longer passing requests through to the external service. Instead, for each request it receives, it is returning a matched response from its cache. The cache is stored in memory, and persisted on disk in the `requets.db` file.

## Next steps

Now you have simulated a service, you can use middleware to simulate network latency or failure, or manipulate data in the responses. Proceed to the **Middleware** section.

Alternatively, you can learn how to use *synthesize mode* to generate responses to requests on the fly. Proceed to the **Creating synthetic services** section.

## Further reading

A detailed step-by-step guide to capturing traffic and creating a simulated service is available here:

[Speeding up your slow dependencies](https://specto.io/blog/speeding-up-your-slow-dependencies.html)

A guide on how to capture and simulate the Meetup API is available here:

[Virtualizing the Meetup API](https://specto.io/blog/hoverfly-meetup-api.html)
