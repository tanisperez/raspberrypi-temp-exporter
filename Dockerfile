FROM golang:1.14.2-buster

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["raspberrypi-temp-exporter"]
