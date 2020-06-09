SHELL = /bin/bash
PACKAGE ?= posology-graphql

#Output purposes
OUTPUT_DIR = $(CURDIR)/output
BIN_OUTPUT_DIR      = $(OUTPUT_DIR)/bin
TEST_OUTPUT_DIR     = $(OUTPUT_DIR)/test
DIRS=$(BIN_OUTPUT_DIR) $(TEST_OUTPUT_DIR)
$(shell mkdir -p $(DIRS))

#Build flags
LDFLAGS ?= "-X 'main.version=$(VERSION)' -X 'main.gitCommit=$(GIT_COMMIT)' -X 'main.application=$(PACKAGE)'"
LINT_FLAGS ?= -E golint -E errcheck --out-format checkstyle
BUILD_FLAGS ?= "./cmd/app"
TEST_FLAGS ?= "-tags=unit"
PACKAGE_EXTENSION ?=  $(shell if [ "$(GOOS)" == windows ]; then echo .exe; fi)
GO111MODULE=on
CGO_ENABLED=0
TARGET ?= final

.SILENT: ; # no need for @
.ONESHELL: ; # recipes execute in same shell
.NOTPARALLEL: ; # wait for this target to finish
.EXPORT_ALL_VARIABLES: ; # send all vars to shell
.PHONY: version build docs test scripts api cmd configs examples

build: ## Build the app
	go build -o $(BIN_OUTPUT_DIR)/$(PACKAGE)$(PACKAGE_EXTENSION) --ldflags=$(LDFLAGS) $(BUILD_FLAGS)

sam-remove-ns:
	sed -i s/ns.://g /media/bertp/Data/SAM/2636-2020-06-03-13-21-08/AMP.xml
	sed -i s/ns.://g /media/bertp/Data/SAM/2636-2020-06-03-13-21-08/REF.xml
	sed -i s/ns.://g /media/bertp/Data/SAM/2636-2020-06-03-13-21-08/VMP.xml
