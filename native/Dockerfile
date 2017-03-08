FROM alpine:3.5

ADD pkg /opt/driver/src/pkg
ADD vendor/cache /opt/driver/src/vendor/cache

RUN mkdir -p /opt/driver/src && \
    #adduser ${BUILD_USER} -u ${BUILD_UID} -D -h /opt/driver/src && \
    apk add --no-cache libc-dev gcc make ruby ruby-dev && \
    cd /opt/driver/src && \
    gem install --no-rdoc --no-ri vendor/cache/* && \
    gem install --no-rdoc --no-ri pkg/* && \
    apk del -v libc-dev gcc make ruby-dev

ENTRYPOINT ["native"]
