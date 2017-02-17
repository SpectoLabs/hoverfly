hoverctl start
hoverctl mode capture
curl --proxy http://hoverfly.example.com:8500 http://ip.jsontest.com
hoverctl mode simulate
curl --proxy http://hoverfly.example.com:8500 http://ip.jsontest.com
hoverctl stop