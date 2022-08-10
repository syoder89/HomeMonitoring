# syntax=docker/dockerfile:1

FROM golang:3.16-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY vmclient ./

RUN go build -o /tasmota-monitor

FROM golang:3.16-alpine
COPY --from=builder /tasmota-monitor /tasmota-monitor
CMD [ "/tasmota-monitor" ]
