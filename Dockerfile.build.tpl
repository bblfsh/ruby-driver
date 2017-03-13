FROM alpine:3.5

ARG RUNTIME_NATIVE_VERSION
ENV RUNTIME_NATIVE_VERSION $RUNTIME_NATIVE_VERSION

RUN mkdir -p /opt/driver/src && \
    adduser ${BUILD_USER} -u ${BUILD_UID} -D -h /opt/driver/src && \
    apk add -v --no-cache libc-dev gcc make git bash \
    ruby=$RUNTIME_NATIVE_VERSION ruby-dev=$RUNTIME_NATIVE_VERSION && \
    gem install -V --no-document rake bundler io-console

WORKDIR /opt/driver/src
