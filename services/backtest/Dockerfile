FROM golang:1.26-alpine AS build

WORKDIR /src

ENV GOPROXY=direct
ENV GOSUMDB=off

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY shared ./shared
COPY services/backtest ./services/backtest

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/backtest-service \
    ./services/backtest/cmd;

FROM alpine:3.21

WORKDIR /app

COPY --from=build /out/backtest-service /app/backtest-service

ENTRYPOINT ["/app/backtest-service"]