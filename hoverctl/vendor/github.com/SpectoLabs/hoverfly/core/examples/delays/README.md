# Delays

This is an example of a file to apply delays within Hoverfly. This example will add a 2 second delay to responses to github.com using any HTTP method. It will also apply a 2 second delay to any response to a POST request.

To import this file into Hoverfly, you can run the following command:
```
hoverctl delays <path to file>
```

To see all delays configured in Hoverfly, use the following command:
```
hoverctl delays
```
To find out more, please check the documentation regarding [simulating service latency](https://spectolabs.gitbooks.io/hoverfly/content/usage/simulating_service_latency.html)
