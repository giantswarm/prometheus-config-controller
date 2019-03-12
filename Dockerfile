FROM alpine:3.8

ADD ./prometheus-config-controller /prometheus-config-controller

ENTRYPOINT ["/prometheus-config-controller"]
