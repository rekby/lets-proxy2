#!/bin/bash
set -ev

GOPATH=$(go env GOPATH)
git clone https://github.com/letsencrypt/boulder/ $GOPATH/src/github.com/letsencrypt/boulder
cd $GOPATH/src/github.com/letsencrypt/boulder

docker-compose build

docker-compose run --use-aliases -e FAKE_DNS=172.17.0.1 --service-ports boulder ./start.py

echo -n "Wait for bounder start listen "
date
while ! curl -q http://localhost:4000 >/dev/null 2>&1; do
    sleep 1
done

