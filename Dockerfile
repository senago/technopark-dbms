FROM golang:1.18-bullseye as builder

WORKDIR /app

COPY . ./
RUN go build -o main cmd/main.go

FROM ubuntu:bionic as exec

WORKDIR /cmd

RUN mkdir /cmd/configs
VOLUME ["/cmd/configs"]

ARG CONFIG
ENV CONFIG=${CONFIG}

COPY --from=builder /app/main ./
COPY --from=builder /app/${CONFIG} ./configs

ENTRYPOINT ["/cmd/main"]
