FROM golang:1.26-alpine AS build

WORKDIR /src

ENV GOPROXY=direct
ENV GOSUMDB=off

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY shared ./shared
COPY services/order-execution ./services/order-execution

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/order-execution-service \
    ./services/order-execution/cmd;

FROM alpine:3.21

WORKDIR /app

COPY --from=build /out/order-execution-service /app/order-execution-service

ENTRYPOINT ["/app/order-execution-service"]