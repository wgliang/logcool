# Makefile for logcool
#
# Author: wgliang
# 	site: wgliang@github.io
# 	
# Targets:
# 	all: Builds the code
# 	build: Builds the code
# 	fmt: Formats the source files
# 	clean: cleans the code
# 	install: Installs the code to the GOPATH
# 	iref: Installs referenced projects
#	test: Runs the tests
#	
#  Blog post on it: http://joneisen.me/post/25503842796
#

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build -v
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test -v 
GOVET=$(GOCMD) vet
GOFMT=gofmt -w

# Package lists
TOPLEVEL_PKG := logcool

# List building
TEST_LIST =  $(filter/...) $(input/...) $(cmd/...) $(output/...)

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST)


$(BUILD_LIST): %_build:
	$(GOBUILD) $(TOPLEVEL_PKG)/$*
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) $(TOPLEVEL_PKG)/$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) $(TOPLEVEL_PKG)/$*
$(TEST_LIST): %_test:
	$(GOTEST) $(TOPLEVEL_PKG)/$*
$(FMT_TEST): %_fmt:
	$(GOFMT) ./$*