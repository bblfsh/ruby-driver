FROM alpine:3.5

ARG RUNTIME_NATIVE_VERSION
ENV RUNTIME_NATIVE_VERSION $RUNTIME_NATIVE_VERSION

ADD build /opt/driver/src

RUN apk add --no-cache libc-dev gcc make ruby="$RUNTIME_NATIVE_VERSION" ruby-dev="$RUNTIME_NATIVE_VERSION" && \
    cd /opt/driver/src && \
    gem install -V --no-document --local dependencies/*  && \
    gem install -V --no-document --bindir /opt/driver/bin --local pkg/* && \
    apk del -v libc-dev gcc make ruby-dev && \
    rm -r pkg dependencies

CMD /opt/driver/bin/driver
