#!/bin/bash
set -ev

GOPATH=$(go env GOPATH)
git clone https://github.com/letsencrypt/boulder/ $GOPATH/src/github.com/letsencrypt/boulder
cd $GOPATH/src/github.com/letsencrypt/boulder

sed -i -e 's/FAKE_DNS.*/FAKE_DNS: 172.17.0.1/' docker-compose.yml # Fake dns to docker host

sed -i -e 's/TRAVIS_GO_VERSION/TRAVIS_GO_VERSION_OFF/' docker-compose.yml # always build boulder with default go version

docker-compose up -d

echo -n "Wait for bounder start listen "
date

while ! curl -q http://localhost:4000 >/dev/null 2>&1; do
    sleep 1
done

date
