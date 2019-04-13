#!/bin/bash
go version

go env

docker --version

ip -4 addr

echo "Prepare boulder"
bash tests/run_docker_boulder.sh
