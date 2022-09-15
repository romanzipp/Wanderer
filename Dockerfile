# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY . ./
RUN go mod download
RUN go build -o /wanderer

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /wanderer /wanderer

EXPOSE 8080

USER nonroot:nonroot

CMD ["/wanderer"]
