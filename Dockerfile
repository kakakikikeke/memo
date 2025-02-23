FROM golang:1.24.0-alpine3.21

ADD ./ $GOPATH/src/github.com/kakakikikeke/memo
WORKDIR $GOPATH/src/github.com/kakakikikeke/memo

RUN go mod tidy
RUN go build

CMD ["./memo"]