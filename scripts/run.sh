#!/bin/bash

# Build the project
go build -o fragmentdb

# Run three nodes with different configurations
./fragmentdb -config config1.json &
./fragmentdb -config config2.json &
./fragmentdb -config config3.json &
