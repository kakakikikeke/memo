FROM golang:1.26.4-alpine3.24

ADD ./ $GOPATH/src/github.com/kakakikikeke/memo
WORKDIR $GOPATH/src/github.com/kakakikikeke/memo

RUN go mod tidy
RUN go build

CMD ["./memo"]