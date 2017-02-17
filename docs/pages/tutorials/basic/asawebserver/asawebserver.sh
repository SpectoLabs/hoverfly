#!/bin/bash

hoverctl start
hoverctl mode capture
curl --proxy http://localhost:8500 http://echo.jsontest.com/a/b
hoverctl export simulation.json
hoverctl stop

hoverctl start webserver
hoverctl import simulation.json
curl http://localhost:8500/a/b
hoverctl stop
