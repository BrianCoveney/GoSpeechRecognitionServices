FROM golang:alpine

MAINTAINER brian.coveney

WORKDIR /go
COPY . .

RUN go get -d -v ./...
RUN go build frontend-service.go

EXPOSE 443
ENTRYPOINT ["/go/frontend-service"]
