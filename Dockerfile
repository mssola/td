FROM golang:1.6-alpine
MAINTAINER Miquel Sabaté Solà <mikisabate@gmail.com>

COPY . /go/src/github.com/mssola/td
RUN go install -ldflags="-s -w" github.com/mssola/td \
      && cp /go/src/github.com/mssola/td/scripts/entrypoint.sh / \
      && rm -rf /go/src \
      && apk add --update vim the_silver_searcher \
      && rm -rf /tmp/* /var/cache/apk/*

ENTRYPOINT ["/entrypoint.sh"]
