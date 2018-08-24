FROM instrumentisto/dep

RUN apt-get update

RUN curl -sSL https://get.docker.com/ | sh

WORKDIR /go/src/Peripli/

RUN git clone https://github.com/Peripli/service-manager
RUN git clone https://github.com/Peripli/service-broker-proxy
RUN git clone https://github.com/Peripli/service-broker-proxy-cf
RUN git clone https://github.com/Peripli/service-broker-proxy-k8s

ENTRYPOINT [ "bash" ]
