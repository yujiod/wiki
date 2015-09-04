FROM golang:1.5

MAINTAINER Yuji ODA

ADD . $GOPATH/src/github.com/yujiod/wiki

RUN go get github.com/revel/cmd/revel
RUN revel build github.com/yujiod/wiki /usr/local/wiki

ENV DB_DRIVER sqlite3
ENV DB_SOURCE ./wiki.db

EXPOSE 9000

WORKDIR /usr/local/wiki
CMD []
ENTRYPOINT ["./run.sh"]
