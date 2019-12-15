#!/bin/bash
env GOOS=linux GOARCH=amd64 go build -v -o bin/webserver
docker build -t consul-example-webserver .
docker tag consul-example-webserver zacsketches/consul-example-webserver
docker push zacsketches/consul-example-webserver
