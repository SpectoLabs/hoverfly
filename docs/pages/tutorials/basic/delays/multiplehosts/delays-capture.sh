hoverctl start
hoverctl mode capture
curl --proxy localhost:8500 http://api.ipify.org
curl --proxy localhost:8500 http://httpbin.org
hoverctl export simulation.json
hoverctl stop
