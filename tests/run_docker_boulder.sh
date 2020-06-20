#!/bin/bash
set -ev

GOPATH=$(go env GOPATH)

LAST_BOULDER_RELEASE_TAG=$(git ls-remote -t https://github.com/letsencrypt/boulder/ | awk '{print $2}' | tr -d '^{}' | sort | tail -n 1 | cut -d '/' -f 3)

echo git checkout boulder release: $LAST_BOULDER_RELEASE_TAG
git clone "--branch=$LAST_BOULDER_RELEASE_TAG" --depth=1  https://github.com/letsencrypt/boulder/ $GOPATH/src/github.com/letsencrypt/boulder
cd $GOPATH/src/github.com/letsencrypt/boulder

sed -i -e 's/FAKE_DNS=.*/FAKE_DNS=172.17.0.1/' docker-compose.yml # Fake dns to docker host

sed -i -e 's/TRAVIS_GO_VERSION/TRAVIS_GO_VERSION_OFF/' docker-compose.yml # always build boulder with default go version

# Set small rate limit windows - for comfort manual test runs
sed -i -e 's/window:.*/window: 1m/' test/rate-limit-policies.yml

docker-compose up -d

echo -n "Wait for boulder start listen "
date

while ! curl -q http://localhost:4001 >/dev/null 2>&1; do
    sleep 1
done

date
