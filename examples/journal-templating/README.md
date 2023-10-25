
# Journal templating

This new templating function allow you to query data from the journal to generate dynamic response.

Use the provided Hoverfly binary, run it in capture mode (which allows request to be passed through and recorded in the journal),
and also import the simulation `journal-templating.json`.

```commandline
 hoverfly -capture -import journal-template.json -journal-indexing-key Request.QueryParam.id 
 curl --proxy localhost:8500 'time.jsontest.com?id=123'
```

The request is now recorded, and we can switch Hoverfly to simulate mode, and then this time we can call the stateful endpoint which will trigger the simulation:
```commandline
 hoverctl mode simulate
 curl --proxy localhost:8500 'time.jsontest.com/stateful'
```

You will see that the stateful endpoint response is the time field from the previous response for `localhost:8500 'time.jsontest.com?id=123`
That's because the simulation uses this templating function:

```json

    "response": {
      "status": 200,
      "body": "{{ journal 'Request.QueryParam.id' '123' 'Response' 'jsonpath' '$.time' }}",
      "encodedBody": false,
      "templated" : true
   }

```

What it does is that it look up the journal with this index key `Request.QueryParam.id`, 
so if any of the previous request has a request query param with key `id` and value `123`, 
then it will return the response body of that request, and then apply the jsonpath expression `$.time` to extract the time field from the response body.

The index key `Request.QueryParam.id` is just an example, you can specify other request data as index to look up using the same syntax as templating request data:
https://docs.hoverfly.io/en/latest/pages/keyconcepts/templating/templating.html#getting-data-from-the-request

eg. if you want to look up the journal by the second request path variable, you can use this index key: `Request.Path.[1]`

```commandline
 hoverfly -journal-indexing-key Request.Path.[1]
```

You can also specify multiple index keys like this: 

```commandline
 hoverfly -journal-indexing-key Request.Path.[1] -journal-indexing-key Request.Header.X-Header-Id 
```

In which case we will create multiple indices for the journal based on the given keys. It allows you to lookup all the past request/response using those keys.