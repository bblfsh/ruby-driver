FROM alpine:3.5

RUN mkdir -p /opt/driver/src && \
    adduser ${BUILD_USER} -u ${BUILD_UID} -D -h /opt/driver/src && \

RUN apk add --no-cache libc-dev gcc make ruby ruby-dev git
RUN gem install bundler io-console

WORKDIR /opt/driver/src 
