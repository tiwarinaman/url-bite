#!/bin/bash
echo "Building the application..."
go build -o url-bite .

echo "Starting the application..."
./url-bite