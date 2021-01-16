FROM golang:1.15-buster

WORKDIR /root/data
COPY . /root/data/

RUN go build .

EXPOSE 8080

ENTRYPOINT ["/root/data/client_api"]

