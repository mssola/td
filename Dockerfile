FROM golang:1.6-alpine
MAINTAINER Miquel Sabaté Solà <mikisabate@gmail.com>

COPY . /go/src/github.com/mssola/td
RUN go install github.com/mssola/td && rm -rf /go/src

ENTRYPOINT ["td"]
