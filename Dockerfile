# This file can be used directly with Docker.
#
# Prerequisites:
#   go mod vendor
#   bblfsh-sdk release
#
# However, the preferred way is:
#   go run ./build.go driver:tag
#
# This will regenerate all necessary files before building the driver.

#==============================
# Stage 1: Native Driver Build
#==============================
FROM ruby:2.4-alpine3.7 as native

# install build dependencies
RUN apk add --no-cache make libc-dev gcc
RUN gem install -V --no-document bundler -v 1.17.3 && gem install -V --no-document rake io-console parser json


ADD native /native
WORKDIR /native

# build native driver
RUN export BUNDLE_IGNORE_CONFIG=1 && bundle install --path vendor/bundle --without development --verbose
RUN rake build --trace
RUN gem install  -V --no-document --local --ignore-dependencies --install-dir ./gems vendor/bundle/ruby/2.4.0/cache/*
RUN gem install  -V --no-document --local --ignore-dependencies --install-dir ./gems --bindir ./build ./pkg/*


#================================
# Stage 1.1: Native Driver Tests
#================================
FROM native as native_test

# install test dependencies
RUN export BUNDLE_IGNORE_CONFIG=1 && bundle install --path vendor/bundle --verbose

# run native driver tests
RUN export GEM_PATH=./vendor/bundle/ruby/2.4.0 && rake test --trace


#=================================
# Stage 2: Go Driver Server Build
#=================================
FROM golang:1.12-alpine as driver

ENV DRIVER_REPO=github.com/bblfsh/ruby-driver
ENV DRIVER_REPO_PATH=/go/src/$DRIVER_REPO

ADD go.* $DRIVER_REPO_PATH/
ADD vendor $DRIVER_REPO_PATH/vendor
ADD driver $DRIVER_REPO_PATH/driver

WORKDIR $DRIVER_REPO_PATH/

ENV GO111MODULE=on GOFLAGS=-mod=vendor

# workaround for https://github.com/golang/go/issues/28065
ENV CGO_ENABLED=0

# build server binary
RUN go build -o /tmp/driver ./driver/main.go
# build tests
RUN go test -c -o /tmp/fixtures.test ./driver/fixtures/

#=======================
# Stage 3: Driver Build
#=======================
FROM ruby:2.4-alpine3.7



LABEL maintainer="source{d}" \
      bblfsh.language="ruby"

WORKDIR /opt/driver

# copy static files from driver source directory
ADD ./native/bin/native.sh ./bin/native
ADD ./native/exe/native ./bin/native.rb
ADD ./native/lib ./bin/lib


# copy build artifacts for native driver
COPY --from=native /native/gems ./bin/gems


# copy driver server binary
COPY --from=driver /tmp/driver ./bin/

# copy tests binary
COPY --from=driver /tmp/fixtures.test ./bin/
# move stuff to make tests work
RUN ln -s /opt/driver ../build
VOLUME /opt/fixtures

# copy driver manifest and static files
ADD .manifest.release.toml ./etc/manifest.toml

ENTRYPOINT ["/opt/driver/bin/driver"]