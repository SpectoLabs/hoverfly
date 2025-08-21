#!/bin/bash


# Rebuild Hoverfly and Hoverctl
make build

# Ensure executables have execution rights
chmod +x target/hoverfly
chmod +x target/hoverctl

# Stop any running Hoverfly instance
./target/hoverctl stop

# Start Hoverfly in webserver mode
./target/hoverctl start webserver

# Add simulation and data source
./target/hoverctl simulation add ~/dev/hoverfly-simulations/product-api-simulation.json
./target/hoverctl templating-data-source set --name products --filePath ~/dev/hoverfly-simulations/products.csv