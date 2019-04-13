#!/bin/bash
go version

go env

docker --version

echo "Prepare boulder"
bash tests/run_docker_boulder.sh
