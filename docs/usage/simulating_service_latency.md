# Simulating service latency
Once you have created a simulated service by capturing traffic between your application and an external service, you may wish to make the simulation more "realistic" by applying latency to the responses returned by Hoverfly.

This could be done using middleware (See the **Using middleware** section). However, if you do not want to go to the effort of writing a middleware script, you can use a JSON file to apply a set of fixed response delays to the Hoverfly simulation.

This method is useful if Hoverfly is being used in a load test to simulate an external service, and there is a requirement to simulate external service latency. Under high load, the overhead of executing middleware scripts will impact the performance of Hoverfly, making the middleware approach to adding latency unsuitable.  

## Set up
To simulate service latency, you will need to have created a simulation by capturing traffic (see the **Capturing traffic** and **Simulating services** sections).

## Simulate latency

### JSON configuration
To simulate latency, Hoverfly can be configured to apply a delay to individual hosts or API endpoints in the simulation using a JSON configuration file. This is done using a regular expression to match against the URL, a delay value in milliseconds, and an optional HTTP method value.

For example, to apply a delay of 2 seconds to all hosts in the simulation:  

    {
      "data": [
        {
          "urlPattern": ".",
          "delay": 2000
        }
      ]
    }

To apply a delay of 1 second to `1.myhost.com` and a delay of 2 seconds to `2.myhost.com`:

    {
      "data": [
        {
          "urlPattern": "1\\.myhost\\.com",
          "delay": 1000
        },
        {
          "urlPattern": "2\\.myhost\\.com",
          "delay": 2000
        }
      ]
    }

It is also possible to apply delays to specific resources and endpoints in your API. In the following example, a delay of 1 second is applied to all endpoints of resource A. For resource B, a delay of 1 second is applied to the GET endpoint, but a different delay of 2 seconds is applied to the POST endpoint:

    {
      "data": [
        {
          "urlPattern": "myhost\\.com\\/A",
          "delay": 1000
        },
        {
          "urlPattern": "myhost\\.com\\/B",
          "delay": 1000,
          "httpMethod": "GET"
        },
        {
          "urlPattern": "myhost\\.com\\/B",
          "delay": 2000,
          "httpMethod": "POST"
        }
      ]
    }

The **delays will be matched in the order that they appear in the JSON configuration file**. In the following example, `"urlPattern":"."` matches all hosts, overriding `"urlPattern": "1\\.myhost\\.com"` and all subsequent matches, applying a 3 second delay to all responses:   

    {
      "data": [
        {
          "urlPattern": ".",
          "delay": 3000
        },
        {
          "urlPattern": "1\\.myhost\\.com",
          "delay": 1000
        }
      ]
    }

### Applying the configuration

The configuration can be applied using hoverctl.

1. To apply delays:
        hoverctl delays path/to/my_delays.json
2. To view the delays which have been applied:
        hoverctl delays

Alternatively, the configuration can be applied using the Hoverfly API directly:

1. To apply delays:
        curl -H "Content-Type application/json" -X PUT -d '{"data":[{"urlPattern":"1\\.myhost\\.com","delay":1000},{"urlPattern":"2\\.myhost\\.com","delay":2000}]}' http://${HOVERFLY_HOST}:8888/api/delays

2. To view the delays which have been applied

        curl http://${HOVERFLY_HOST}:8888/api/delays
