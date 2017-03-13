SOURCEDIR=.
GEMDIR=pkg
VENDORDIR=vendor
DRIVER_VERSION=beta-0.0.9
DOCKERFILE=Dockerfile
DOCKER_IMAGE_NAME=bblfsh/ruby-driver
DOCKER_IMAGE_BUILD_NAME=bblfsh/ruby-driver-build

all: bundle-gem build-gem docker-build

bundle-gem:
	bundle install --no-deployment --path vendor/bundle --verbose;

build-gem:
	rake build --trace

docker-build: $(DOCKERFILE)
	docker build -f ${DOCKERFILE} -t ${DOCKER_IMAGE_NAME} ${SOURCEDIR}

.PHONY:
clean:
	if [ -d ${GEMDIR} ] ; then rm -r ${GEMDIR} ; fi
	if [ -d ${VENDORDIR} ] ; then rm -r ${VENDORDIR}; fi
	if [ -f Gemfile.lock ] ; then rm Gemfile.lock ; fi
