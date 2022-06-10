FROM golang:1.18-alpine as builder

WORKDIR /app

COPY . ./
RUN go build -o main cmd/main.go

FROM alpine:latest as exec

WORKDIR /cmd

RUN mkdir /cmd/configs
VOLUME ["/cmd/configs"]

ARG CONFIG
ENV CONFIG=${CONFIG}

COPY --from=builder /app/main ./
COPY --from=builder /app/${CONFIG} ./configs

ENTRYPOINT ["/cmd/main"]
