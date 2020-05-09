FROM golang:alpine3.11

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["raspberrypi-temp-exporter"]