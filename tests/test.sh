#!/bin/bash
set -ev

CURRENT_DIR=$(realpath "$0")
CURRENT_DIR=$(dirname "$CURRENT_DIR")

cd "$CURRENT_DIR/.."

docker-compose up --abort-on-container-exit
