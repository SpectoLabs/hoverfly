# Steps for invoking post serve action

- Make sure requests package is there. If not, then install it via pip. 
```shell
    pip install requests
```
- Start hoverfly by passing post-serve-action
 ```shell
    ./target/hoverfly -post-serve-action "outbound-http python3 examples/postserveaction/outboundhttpaction.py 1000"
```
- Verify post serve action is registered by making GET call - <host>/api/v2/hoverfly/post-serve-actions

- Once registered, you need to change simulation file and include post serve action name in which you want to invoke it.

Example
```json
{
  "data": {
    "pairs": [
      {
        "request": {
          ...
          "destination": [
            {
              "matcher": "exact",
              "value": "helloworld-test.com"
            }
          ]
          ...
        },
        "response": {
          "status": 200,
          "postServeAction": "outbound-http",
          "body": "Hello World",
          "encodedBody": false,
          ...
        }
      }
    ],
    ...
  },
  "meta": {
    "schemaVersion": "v5.2",
    "hoverflyVersion": "v1.6.0",
    "timeExported": "2023-09-02T13:10:04+05:30"
  }
}

```