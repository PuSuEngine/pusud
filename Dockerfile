FROM golang:alpine
RUN apk add --no-cache git
ADD . /go/src/app
WORKDIR /go/src/app

RUN go-wrapper download

RUN go build -v -o /go/bin/pusud pusud.go

ENTRYPOINT /go/bin/pusud

EXPOSE 55000
