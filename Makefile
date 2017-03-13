RUBY_MAKE_CMD=rake
RUBY_DEP_PACK_CMD=bundle
RUBY_GEM_CMD=gem

test-native-internal:
	if [ -f .bundle/config ] ; then rm .bundle/config ; fi
	cd native; \
	$(RUBY_DEP_PACK_CMD) install --path vendor/bundle --verbose; \
	export GEM_PATH=./vendor/bundle/ruby/2.3.0 && $(RUBY_MAKE_CMD) test --trace;

build-native-internal:
	if [ -f .bundle/config ] ; then rm .bundle/config ; fi
	cd native; \
	$(RUBY_DEP_PACK_CMD) install --path vendor/bundle --verbose; \
	$(RUBY_MAKE_CMD) build --trace; \
	cp -r pkg $(BUILD_PATH); \
	mkdir -p $(BUILD_PATH)/dependencies; \
	cp -r vendor/bundle/ruby/2.3.0/cache/json* $(BUILD_PATH)/dependencies;

include .sdk/Makefile
