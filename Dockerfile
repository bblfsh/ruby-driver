FROM alpine:latest

RUN mkdir -p /opt/driver/src
#    && adduser ${BUILD_USER} -u ${BUILD_UID} -D -h /opt/driver/src

RUN apk add --no-cache libc-dev gcc make ruby ruby-dev git
#RUN gem install --no-rdoc --no-ri bundler io-console

ADD pkg /opt/driver/src/pkg
ADD vendor/cache /opt/driver/src/vendor/cache

WORKDIR /opt/driver/src

#RUN bundle install --local
# RUN rake install

# RUN rm -r /opt/driver/src

#RUN gem uninstall bundler io-console

RUN gem install --no-rdoc --no-ri vendor/cache/*
RUN gem install --no-rdoc --no-ri pkg/*
RUN apk del -v libc-dev gcc make git

ENTRYPOINT ["ruby-driver", "docker-driver"]
