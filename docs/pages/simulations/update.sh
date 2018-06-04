#! /bin/sh

echo "You have the latest Hoverfly, right?"
hoverctl start
hoverctl import pages/simulations/all-matchers-simulation.json
hoverctl export pages/simulations/all-matchers-simulation.json
hoverctl import pages/simulations/basic-delay-simulation.json
hoverctl export pages/simulations/basic-delay-simulation.json
hoverctl import pages/simulations/basic-encoded-simulation.json
hoverctl export pages/simulations/basic-encoded-simulation.json
hoverctl import pages/simulations/basic-simulation.json
hoverctl export pages/simulations/basic-simulation.json
hoverctl import pages/simulations/get-method-delay-simulation.json
hoverctl export pages/simulations/get-method-delay-simulation.json
hoverctl import pages/simulations/multiple-hosts-delay-simulation.json
hoverctl export pages/simulations/multiple-hosts-delay-simulation.json
hoverctl import pages/simulations/multiple-locations-delay-simulation.json
hoverctl export pages/simulations/multiple-locations-delay-simulation.json
hoverctl stop
sed -i 's/\t/   /g' pages/simulations/*.json