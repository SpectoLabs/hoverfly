hoverctl start
hoverctl import simulation.json
curl --proxy localhost:8500 http://api.ipify.org
curl --proxy localhost:8500 http://httpbin.org
hoverctl stop
