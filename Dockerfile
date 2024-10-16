# Container image that runs your code
FROM golang:1.23-alpine

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -o /go/bin/app

ENTRYPOINT ["/go/bin/app"]