# syntax=docker/dockerfile:1

FROM golang:1.21 AS build-stage
WORKDIR /app

COPY go.mod ./
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /cowsay-http

CMD ["/cowsay-http"]

# Deploy the application binary into a lean image
FROM debian:latest AS build-release-stage

ENV PATH="${PATH}:/usr/games/"

RUN <<EOF
  apt-get update
  apt-get install -y fortune cowsay cowsay-off
EOF

WORKDIR /

COPY --from=build-stage /cowsay-http /cowsay-http
COPY pages/ ./pages/

ENTRYPOINT ["/cowsay-http"]
