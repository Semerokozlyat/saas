#Build
FROM golang:1.12-alpine as build_location
#FROM golang:latest as build_location
ARG PACKAGE_NAME=saas
ARG BUILD_VERSION=0.1

RUN apk add --update git gcc musl-dev
ADD . /go/src/${PACKAGE_NAME}
WORKDIR /go/src/${PACKAGE_NAME}

RUN CGO_ENABLED=1 GOOS=linux go build -a \
# -installsuffix cgo \
 -gcflags '-N -l' \
 -ldflags '-w -X main.Version=${BUILD_VERSION}' -o saas ./app/

# Install service
FROM alpine:latest

# Install Chrome (based on https://github.com/Zenika/alpine-chrome/blob/master/Dockerfile)
ARG BUILD_DATE
ARG VCS_REF

# Installs latest Chromium package.

RUN echo @edge http://nl.alpinelinux.org/alpine/edge/community > /etc/apk/repositories \
    && echo @edge http://nl.alpinelinux.org/alpine/edge/main >> /etc/apk/repositories \
    && apk add --no-cache \
    libstdc++@edge \
    chromium@edge \
    harfbuzz@edge \
    nss@edge \
    freetype@edge \
    ttf-freefont@edge \
    && rm -rf /var/cache/* \
    && mkdir /var/cache/apk

# Add Chrome as a user
RUN mkdir -p /usr/src/app \
    && adduser -D chrome \
    && chown -R chrome:chrome /usr/src/app
# Run Chrome as non-privileged
USER chrome
WORKDIR /usr/src/app

ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/lib/chromium/

WORKDIR /opt/saas/${PACKAGE_NAME}
COPY --from=build_location /go/src/${PACKAGE_NAME}/saas /opt/saas/
#EXPOSE 8000

# Autorun chrome headless with no GPU
#ENTRYPOINT ["chromium-browser", "--headless", "--disable-gpu", "--disable-software-rasterizer", "--disable-dev-shm-usage"]
ENTRYPOINT ["/opt/saas/saas"]
