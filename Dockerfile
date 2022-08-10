# syntax=docker/dockerfile:1

FROM golang:1.19.0-alpine3.16 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY vmclient ./

RUN go build -o /tasmota-monitor

FROM golang:1.19.0-alpine3.16
COPY --from=builder /tasmota-monitor /tasmota-monitor
CMD [ "/tasmota-monitor" ]
