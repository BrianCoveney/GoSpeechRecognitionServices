FROM golang:alpine

MAINTAINER brian.coveney

WORKDIR /go
COPY . .

RUN apk add --no-cache git mercurial \
    && go get github.com/random/library/that/I/use/ \
    && apk del git mercurial
RUN go get -d -v ./...
RUN go build frontend-service.go


EXPOSE 443
ENTRYPOINT ["/go/frontend-service"]
