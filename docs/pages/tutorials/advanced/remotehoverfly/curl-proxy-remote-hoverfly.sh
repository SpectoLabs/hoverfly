hoverctl mode -t remote capture
curl --proxy http://hoverfly.example.com:8555 http://ip.jsontest.com
hoverctl mode -t remote simulate
curl --proxy http://hoverfly.example.com:8555 http://ip.jsontest.com