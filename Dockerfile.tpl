FROM bblfsh/ruby-driver-build

ADD build /opt/driver
WORKDIR /opt/driver

RUN gem install -V --no-document --local dependencies/*  && \
    gem install -V --no-document --bindir /opt/driver/bin --local pkg/* && \
    apk del -v libc-dev gcc make ruby-dev && \
    rm -r pkg dependencies

ENTRYPOINT ["/opt/driver/bin/driver"]
