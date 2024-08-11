# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

COPY *.go ./

RUN go build -v -o /prusa_log_processor

FROM alpine:latest

COPY --from=builder /prusa_log_processor .

EXPOSE 10010

ENTRYPOINT ["/prusa_log_processor"]