FROM alpine:3.5

RUN mkdir -p /opt/driver/src && \
    adduser ${BUILD_USER} -u ${BUILD_UID} -D -h /opt/driver/src

RUN apk add  --no-cache make git ca-certificates libc-dev gcc \
        ruby=${RUNTIME_NATIVE_VERSION} ruby-dev=${RUNTIME_NATIVE_VERSION}

RUN gem install -V --no-document rake bundler io-console parser

WORKDIR /opt/driver/src
