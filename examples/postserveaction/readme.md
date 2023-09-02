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