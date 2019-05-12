Need for go generate:
* Install https://github.com/gojuno/minimock to PATH
* Install https://github.com/gobuffalo/packr to PATH

Need for boulder tests:
* docker, see tests/run_docker_boulder.sh for run boulder

Fake DNS - set to IP of devel computer, allowed from docker.

Must bind docker port 4000 to local port 4000 (for integration tests).
If use docker-machine - need ```docker-machine ssh <machine-name> -L 4000:localhost:4000```

