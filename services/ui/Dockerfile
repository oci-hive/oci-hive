FROM docker.io/golang:1.27-alpine

WORKDIR /var/www/html

RUN apk update && \
    apk add --no-cache \
        git \
        vim \
        ca-certificates \
        bash \
        nodejs \
        npm && \
    echo 'alias ll="ls -hlaF"' >> ~/.bashrc && \
    echo 'alias vi="vim"' >> ~/.bashrc
