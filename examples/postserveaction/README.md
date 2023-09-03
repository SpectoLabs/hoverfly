# Steps for invoking post serve action

- Install the python requests library via pip 
```shell
    pip install requests
```
- Start hoverfly by passing a post-serve-action, and a simulation file that invokes post serve action by its name.
 ```shell
    hoverfly -post-serve-action "outbound-http python3 examples/postserveaction/outboundhttpaction.py 1000" -import examples/postserveaction/simulation-with-callback.json
```
- You can verify that the post serve action is registered by using the admin endpoint: http://localhost:8888/api/v2/hoverfly/post-serve-actions

- Proxying a request to http://helloworld-test.com should trigger a callback to http://ip.jsontest.com,
  and the log from the post serve action script should be printed out to the hoverfly logs: 
 ```shell
    INFO[2023-09-03T12:46:25+01:00] Output from post serve action HTTP call invoked from IP Address: 100.197.184.111 
```
