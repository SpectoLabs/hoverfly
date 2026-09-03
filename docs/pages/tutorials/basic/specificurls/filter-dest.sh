hoverctl start
hoverctl destination "ip" 
hoverctl mode capture
curl --proxy http://localhost:8500 http://api.ipify.org
curl --proxy http://localhost:8500 http://httpbin.org
hoverctl logs
hoverctl export simulation.json
hoverctl stop
