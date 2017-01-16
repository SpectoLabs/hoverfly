hoverctl start
hoverctl mode capture
curl --proxy localhost:8500 http://time.jsontest.com
hoverctl mode simulate
hoverctl delays delays.json
curl --proxy localhost:8500 http://time.jsontest.com
hoverctl stop
