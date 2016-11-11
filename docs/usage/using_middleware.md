# Using middleware
Hoverfly middleware can be written in any language. Middleware modules receive a service data JSON string via the standard input (STDIN) and must return a service data JSON string to the standard output (STDOUT). 

The way middleware is applied to the requests and/or the responses depends on the mode:

* Capture: Calls middleware once with an outgoing request object
* Simulate: Calls middleware once with both request and response objects but any modifications will only affect the response
* Synthesize: Calls middleware once with a request object so that the middleware can create a response
* Modify Mode: Calls middleware twice, the first with a request object so that the request can be modified and then a second time with both the request and the response

**In *simulate mode*, middleware is only executed when a request is matched to a stored response**.

To implement more dynamic behaviour, middleware should be combined with Hoverfly's *synthesize mode* (see the **Creating synthetic services** section).

## Simple middleware examples

Please read the **Installation and setup**, **Capturing traffic** and **Simulating services** sections (if you haven't already) before proceeding.

Middleware is set using hoverctl.

    hoverctl middleware "path/to/script"

You can also set the middleware on start up with the Hoverfly binary

    ./hoverfly -middleware "path/to/script"
    
This will start Hoverfly in the default mode (*simulate mode*) and import a middleware script.

The string supplied as middleware can contain commands as well as a path to a file. For example, if you have written middleware in Go:

    hoverctl middleware "go run path/to/file.go"
    
(This will compile the Go file on the fly, which will introduce latency. It would be preferable to pre-compile the Go binary.)

### Python example

This example will change the response code and body in each response, and add 2 seconds of delay (simulating network latency).

1. Ensure that you have captured some traffic with Hoverfly
2. Ensure that you have Python installed
3. Save the following code into a file named `example.py` and make it executable (`chmod +x example.py`):

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

4. Retart Hoverfly in *simulate mode* with the `example.py` script specified as middleware:

       hoverctl middleware "./example.py"

5. Repeat the steps you took to capture the traffic       

You will notice that every response will have the `201` status code, and the body will have been replaced by the string specified in the script. 

There will also be a 2 second delay between each request and the response.

### Javascript example

This example will change the response code and body in each response.

1. Ensure that you have captured some traffic with Hoverfly
2. Ensure that you have NodeJS installed
3. Save the following code into a file named `example.js` and make it executable (`chmod +x example.js`):

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
        
4. Restart Hoverfly in *simulate mode* with the `example.js` script specified as middleware:

          hoverctl middleware "./example.js"
          
5. Repeat the steps you took to capture the traffic
You will notice that every response will have the `201` status code, and the body will have been replaced by the string specified in the script.

## Remote execution of middleware
The examples above show how to execute middleware on the machine on which Hoverfly is running. It is also possible to execute middleware remotely over HTTP.

For example, you could write a piece of middleware using [AWS Lambda](https://docs.aws.amazon.com/lambda/latest/dg/welcome.html) and then provide the URL to Hoverfly as middleware.

    hoverctl middleware "https://1sfefc4.execute-api.eu-west-1.amazonaws.com/remote/middleware"
          

Remote middleare execution works in a similar way to local middleware execution. Hoverfly will send the same JSON data to the middleware, but instead of sending it via `stdin`, it is sent as a POST request to the specified URL. The middleware is then expected to return the same data as it would to `stdout` to the response body.


## Base64 encoding binary data

The above example assumes that the responses do not contain binary data.  If they do then the `encodedBody` field will be set to `true`.  In this circumstance, if you want to mutate the body you must base64 encode and decode it.  

It is also worth bearing in mind that the "Content-Length" header must be set to the length of the unencoded body.

## Further reading

A detailed step-by-step guide to using middleware is available here:

[Modifying traffic on the fly](https://specto.io/blog/service-virtualization-is-so-last-year.html)

More middleware examples are available here:

[Middleware examples](https://github.com/SpectoLabs/hoverfly/tree/master/examples/middleware)

Currently, middleware is not supported by hoverctl.