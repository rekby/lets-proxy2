#!/bin/bash
GOPATH=$(go env GOPATH)
git clone https://github.com/letsencrypt/boulder/ $GOPATH/src/github.com/letsencrypt/boulder
cd $GOPATH/src/github.com/letsencrypt/boulder
docker-compose up -d

while ! curl -q http://localhost:4000 >/dev/null 2>&1; do
    echo -n "Wait for bounder start listen "
    date
    sleep 1
done

