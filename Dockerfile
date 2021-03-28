FROM golang:1
WORKDIR /lets-proxy
COPY . /lets-proxy
RUN cd /lets-proxy/cmd && CGO_ENABLED=0 GOOS=linux go build -mod vendor -o lets-proxy

FROM alpine:latest
COPY cmd/static/default-config.toml /etc/lets-proxy.default.config.toml
COPY --from=0 /lets-proxy/cmd/lets-proxy /lets-proxy
CMD ["/lets-proxy", "--config=/etc/lets-proxy.default.config.toml"]
