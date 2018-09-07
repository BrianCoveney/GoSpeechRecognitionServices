FROM golang:alpine

MAINTAINER brian.coveney

WORKDIR /go
COPY . .

RUN apk add --update --no-cache git
RUN apk add --update --no-cache bzr
RUN go get -d -v ./...
RUN go build frontend-service.go


EXPOSE 443
ENTRYPOINT ["/go/frontend-service"]
