# Request Templates
This is an example of a request template. This request template allows Hoverfly to serve a response to a partially matched request.

In this example, two templates are used for matching. The first will match a GET request to specto.io/virtualized. The second will match any request that includes the header "Match: true".

To import these templates into Hoverfly, you can run the following command:
```
hoverctl templates <path to file>
```
To see all request templates configured in Hoverfly, use the following command:
```
hoverctl templates
```
To find out more, please check the documentation regarding [partial matching](https://spectolabs.gitbooks.io/hoverfly/content/usage/matching_requests.html)
