FROM golang:1.26-alpine AS build

WORKDIR /src

ENV GOPROXY=direct
ENV GOSUMDB=off

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY shared ./shared
COPY services/chart ./services/chart

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/chart-service \
    ./services/chart/cmd;

FROM alpine:3.21

WORKDIR /app

COPY --from=build /out/chart-service /app/chart-service

ENTRYPOINT ["/app/chart-service"]