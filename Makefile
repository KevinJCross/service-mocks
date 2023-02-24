SHELL := /bin/bash
.SHELLFLAGS = -euo pipefail -c
MAKEFLAGS = -s
.SHELLFLAGS := -eu -o pipefail -c ${SHELLFLAGS}
test_dirs=$(shell   find . -name "*_test.go" -exec dirname {} \; |  cut -d/ -f2 | sort | uniq)
export GO111MODULE=on

test_dirs=$(shell   find . -name "*_test.go" -exec dirname {} \; |  cut -d/ -f2 | sort | uniq)

.PHONY: build
build:
	@echo "# building build/aws_lambda/main"
	mkdir -p build/aws_lambda
	GOOS=linux GOARCH=amd64 go build -o build/aws_lambda/main cmd/aws_lambda/main.go

.PHONY: build_tests
build_tests: $(addprefix build_test-,$(test_dirs))

build_test-%:
	@echo " - building '$*' tests"
	@export build_folder=${PWD}/build/tests/$* &&\
	 mkdir -p $${build_folder} &&\
	 cd $* &&\
	 for package in $$(  go list ./... | sed 's|.*/autoscaler/$*|.|' | awk '{ print length, $$0 }' | sort -n -r | cut -d" " -f2- );\
	 do\
	   export test_file=$${build_folder}/$${package}.test;\
	   echo "   - compiling $${package} to $${test_file}";\
	   go test -c -o $${test_file} $${package};\
	 done;

.PHONY: ginkgo_check
ginkgo_check:
	current_version=$(shell ginkgo version | cut -d " " -f 3 | sed -E 's/([0-9]+\.[0-9]+)\..*/\1/');\
	expected_version=$(shell grep "ginkgo"  ".tool-versions" | cut -d " " -f 2 | sed -E 's/([0-9]+\.[0-9]+)\..*/\1/');\
	if [ "$${current_version}" != "$${expected_version}" ]; then \
        echo "ERROR: Expected to have ginkgo version '$${expected_version}.x' but we have $(shell ginkgo version)";\
        exit 1;\
    fi

.PHONY: test
test: ginkgo_check
	@echo "# Running tests"
	@ginkgo run -p -r --race --require-suite --randomize-all --cover ${OPTS}

.PHONY: clean
clean:
	@echo "# cleaning"
	@go clean -cache -testcache
	@rm -rf build

.PHONY: tools
tools:
	which ginkgo > /dev/null || go install github.com/onsi/ginkgo/v2/ginkgo
	which goimports > /dev/null || go install golang.org/x/tools/cmd/goimports

imports: tools
	@echo " - goimports ."
	goimports -w .

.PHONY: lint-actions
lint-actions:
	@echo "- linting GitHub actions"
	go run github.com/rhysd/actionlint/cmd/actionlint@latest

lint: tools lint-actions
	@echo " - linting: ."
	@golangci-lint run --config .golangci.yaml ${OPTS}
