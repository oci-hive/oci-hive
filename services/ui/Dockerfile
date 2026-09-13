FROM docker.io/golang:1.27-alpine

WORKDIR /var/www/html

RUN apk update && \
    apk add --no-cache \
        git \
        ca-certificates \
        bash \
        nodejs \
        npm
