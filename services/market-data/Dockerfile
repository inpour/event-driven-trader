FROM golang:1.26-alpine AS build

WORKDIR /src

ENV GOPROXY=direct
ENV GOSUMDB=off

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY shared ./shared
COPY services/market-data ./services/market-data

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/market-data-service \
    ./services/market-data/cmd;

FROM alpine:3.21

WORKDIR /app

COPY --from=build /out/market-data-service /app/market-data-service

ENTRYPOINT ["/app/market-data-service"]