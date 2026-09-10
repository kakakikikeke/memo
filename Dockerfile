FROM golang:1.26.4-alpine3.24

COPY . $GOPATH/src/github.com/kakakikikeke/memo
WORKDIR $GOPATH/src/github.com/kakakikikeke/memo

RUN addgroup -S memo && adduser -S -G memo memo
RUN go mod tidy
RUN go build

USER memo
CMD ["./memo"]