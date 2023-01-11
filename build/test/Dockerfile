FROM golang:1.19.3-alpine3.16

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ADD . /project/
WORKDIR /project/

RUN go test -v -coverprofile=coverage.out ./gocollection/... \
    && go tool cover -func=coverage.out
