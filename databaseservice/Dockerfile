FROM golang:1.9.3

MAINTAINER brian.coveney

WORKDIR /go
COPY . .

RUN go get -d -v ./...
RUN go build database-service.go

EXPOSE 3001
ENTRYPOINT ["/go/database-service"]
