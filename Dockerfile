# syntax=docker/dockerfile:1

## Build go application
FROM golang:1.19-buster AS build

WORKDIR /app

COPY . ./
RUN go mod download
RUN go build -o wanderer

## Build frontend
FROM node:20 AS build-frontend

WORKDIR /app

COPY . ./
RUN npm install
RUN npm run build

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /app/wanderer /app/wanderer
COPY --from=build /app/views /app/views
COPY --from=build-frontend /app/dist/app.css /app/dist/app.css
COPY --from=build-frontend /app/static /app/static

EXPOSE 8080

USER nonroot:nonroot

CMD ["./wanderer"]
