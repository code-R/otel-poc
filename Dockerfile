FROM golang:1.20-alpine as builder

ENV GOPROXY=https://proxy.golang.org

RUN mkdir -p /src

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

ADD . .

RUN go build -o /bin/app ./cmd

CMD ["app"]
