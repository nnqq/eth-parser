FROM golang:1.20-alpine AS build
WORKDIR /app
COPY / /app
RUN go build -o servicebin cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/servicebin /app
