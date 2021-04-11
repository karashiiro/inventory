FROM golang:1.16-alpine AS build
ENV CGO_ENABLED=0
WORKDIR /src
COPY ./ /src
RUN go mod download
RUN go build -ldflags="-s -w"
FROM alpine:latest
WORKDIR /app
COPY --from=build /src/inventory /app
CMD ["/app/inventory"]
