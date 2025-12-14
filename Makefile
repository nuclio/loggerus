# Copyright 2021 The Nuclio Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	 http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GOPATH ?= $(shell go env GOPATH)
OS_NAME = $(shell uname)
ROOT_DIR := $(shell pwd)

GOLANGCI_LINT_VERSION := 2.7.2
GOLANGCI_LINT_BIN := $(ROOT_DIR)/.bin/golangci-lint
GOLANGCI_LINT_INSTALL_COMMAND := curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(ROOT_DIR)/.bin v$(GOLANGCI_LINT_VERSION)


.PHONY: fmt
fmt: ensure-golangci-linter
	gofmt -s -w .
	$(GOLANGCI_LINT_BIN) run --fix

.PHONY: lint
lint: modules ensure-golangci-linter
	@echo Linting...
	$(GOLANGCI_LINT_BIN) run -v
	@echo Done.

.PHONY: test
test: modules ensure-gopath
	go test -v ./... -short

.PHONY: modules
modules: ensure-gopath
	@echo Getting go modules
	@go mod download

.PHONY: ensure-gopath
ensure-gopath:
ifndef GOPATH
	$(error GOPATH must be set)
endif

.PHONY: ensure-golangci-linter
ensure-golangci-linter:
	@mkdir -p $(ROOT_DIR)/.bin
	@if ! command -v $(GOLANGCI_LINT_BIN) >/dev/null 2>&1; then \
		echo "golangci-lint not found. Installing..."; \
		$(GOLANGCI_LINT_INSTALL_COMMAND); \
	else \
		installed_version=$$($(GOLANGCI_LINT_BIN) version | awk '/version/ {gsub(/^v/, "", $$4); print $$4}'); \
		if [ "$$installed_version" != "$(GOLANGCI_LINT_VERSION)" ]; then \
			echo "Version mismatch ($$installed_version != $(GOLANGCI_LINT_VERSION)). Reinstalling..."; \
			$(GOLANGCI_LINT_INSTALL_COMMAND); \
		fi \
	fi
