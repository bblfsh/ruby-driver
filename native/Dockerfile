FROM alpine:3.5

ADD pkg /opt/driver/src/pkg
ADD vendor/bundle/ruby/2.3.0/cache/ /opt/driver/src/cache

RUN mkdir -p /opt/driver/src && \
    apk add -v --no-cache libc-dev gcc make ruby ruby-dev && \
    cd /opt/driver/src && \
    gem install -V --no-document --local cache/* && \
    gem install -V --no-document --bindir /opt/driver/bin --local pkg/* && \
    apk del -v libc-dev gcc make ruby-dev && \
    rm -r /opt/driver/src/*

ENTRYPOINT ["/opt/driver/bin/native"]
