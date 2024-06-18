FROM golang:1.22.4-alpine3.20

ADD ./ $GOPATH/src/github.com/kakakikikeke/memo
WORKDIR $GOPATH/src/github.com/kakakikikeke/memo

RUN go mod tidy
RUN go build

CMD ["./memo"]