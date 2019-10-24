#!/usr/bin/env bash

docker run -d -p 9091:8080 -e MESSAGE="port 9091" test
docker run -d -p 9092:8080 -e MESSAGE="port 9092" test
