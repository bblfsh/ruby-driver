FROM alpine:latest

RUN mkdir -p /opt/driver/src
#    && adduser ${BUILD_USER} -u ${BUILD_UID} -D -h /opt/driver/src

RUN apk add --no-cache libc-dev gcc make ruby ruby-dev

ADD pkg /opt/driver/src/pkg
ADD vendor/cache /opt/driver/src/vendor/cache

WORKDIR /opt/driver/src

RUN gem install --no-rdoc --no-ri vendor/cache/*
RUN gem install --no-rdoc --no-ri pkg/*
RUN apk del -v libc-dev gcc make ruby-dev

ENTRYPOINT ["native"]
