hoverctl start
hoverctl mode simulate
hoverctl import simulationimport.json
curl --proxy http://localhost:8500 http://echo.jsontest.com/foo/QUX/bar/spam
hoverctl stop
