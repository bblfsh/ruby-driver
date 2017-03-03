FROM alpine

ADD . /ruby-driver
RUN apk add --no-cache libc-dev gcc make ruby ruby-dev git
RUN gem install --no-rdoc --no-ri bundler io-console
WORKDIR ruby-driver
RUN bundle update
RUN rake install

WORKDIR /
RUN rm -r /ruby-driver
RUN apk del -v libc-dev gcc make git
RUN gem uninstall bundler io-console

ENTRYPOINT ["ruby-driver", "docker-driver"]
