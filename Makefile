-include .sdk/Makefile

$(if $(filter true,$(sdkloaded)),,$(error You must install bblfsh-sdk))

RUBY_MAKE_CMD=rake
RUBY_DEP_PACK_CMD=bundle
RUBY_GEM_CMD=gem

test-native-internal:
	if [ -f native/.bundle/config ] ; then rm native/.bundle/config ; fi
	cd native; \
	export BUNDLE_IGNORE_CONFIG=1 && $(RUBY_DEP_PACK_CMD) install --path vendor/bundle --verbose; \
	export GEM_PATH=./vendor/bundle/ruby/2.4.0 && $(RUBY_MAKE_CMD) test --trace;

build-native-internal:
	if [ -f native/.bundle/config ] ; then rm native/.bundle/config ; fi
	cd native; \
	export BUNDLE_IGNORE_CONFIG=1 && $(RUBY_DEP_PACK_CMD) install --path vendor/bundle --without development --verbose; \
	$(RUBY_MAKE_CMD) build --trace; \
	cp -r pkg $(BUILD_PATH); \
	mkdir -p $(BUILD_PATH)/dependencies; \
	cp -r vendor/bundle/ruby/2.4.0/cache/* $(BUILD_PATH)/dependencies;
