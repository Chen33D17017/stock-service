FROM golang:1.15-alpine

ENV GO111MODULE=on

RUN mkdir /go-stock
ADD . /go-stock

COPY go.mod /go-stock
COPY go.sum /go-stock

WORKDIR /go-stock
RUN go clean --modcache
RUN go mod download

RUN go build -o main .
CMD ["/go-stock/main"]