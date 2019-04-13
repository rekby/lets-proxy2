#!/bin/bash
go version

go env

docker --version

echo "Prepare boulder"
tests/run_docker_boulder.sh
