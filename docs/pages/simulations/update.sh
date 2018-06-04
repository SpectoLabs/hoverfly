#! /bin/sh

echo "You have the latest Hoverfly, right?"
hoverctl start
hoverctl import all-matchers-simulation.json
hoverctl export all-matchers-simulation.json
hoverctl import basic-delay-simulation.json
hoverctl export basic-delay-simulation.json
hoverctl import basic-encoded-simulation.json
hoverctl export basic-encoded-simulation.json
hoverctl import basic-simulation.json
hoverctl export basic-simulation.json
hoverctl import get-method-delay-simulation.json
hoverctl export get-method-delay-simulation.json
hoverctl import multiple-hosts-delay-simulation.json
hoverctl export multiple-hosts-delay-simulation.json
hoverctl import multiple-locations-delay-simulation.json
hoverctl export multiple-locations-delay-simulation.json
hoverctl stop