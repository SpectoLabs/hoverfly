hoverctl start
hoverctl import simulation.json
curl --proxy localhost:8500 https://api.ipify.org
curl --proxy localhost:8500 https://httpbin.org
hoverctl stop
