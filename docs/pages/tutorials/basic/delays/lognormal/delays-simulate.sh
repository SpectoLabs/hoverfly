hoverctl start
hoverctl import simulation.json
curl --proxy localhost:8500 http://echo.jsontest.com/b/c
curl -X POST --proxy localhost:8500 http://echo.jsontest.com/b/c
hoverctl stop
