FROM golang:1.15-buster

WORKDIR /root/data
COPY . /root/data/

RUN go build .

ENTRYPOINT ["/root/data/port_domain_service"]

