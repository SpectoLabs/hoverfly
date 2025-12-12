#!/bin/bash

# Stop any running Hoverfly instance
./target/hoverctl stop
./target/hoverctl stop -t local

# Rebuild Hoverfly and Hoverctl
make build

# Ensure executables have execution rights
chmod +x target/hoverfly
chmod +x target/hoverctl


# Start Hoverfly in webserver mode
./target/hoverctl start webserver --log-level debug

# Add simulation 
./target/hoverctl simulation add ~/dev/hoverfly-debug-stuff/jwt-jsonpath.json
