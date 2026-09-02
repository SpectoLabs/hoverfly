hoverctl start
hoverctl mode capture
curl --proxy localhost:8500 https://api.ipify.org
curl --proxy localhost:8500 https://httpbin.org
hoverctl export simulation.json
hoverctl stop
