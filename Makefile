-include .sdk/Makefile

RUBY_TEST_COMMAND=rake

$(if $(filter true,$(sdkloaded)),,$(error You must install bblfsh-sdk))

test-native:
	cd native; \
	$(RUBY_TEST_COMMAND) test

build-native:
	cd native; \
	echo "not implemented"
	echo -e "#!/bin/bash\necho 'not implemented'" > $(BUILD_PATH)/native
	chmod +x $(BUILD_PATH)/native
