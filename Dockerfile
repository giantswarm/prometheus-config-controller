FROM scratch

ADD ./prometheus-config-controller /prometheus-config-controller

ENTRYPOINT ["/prometheus-config-controller"]
