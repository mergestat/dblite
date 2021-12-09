.PHONY: clean update vet test lint lint-ci test-cover bench

# default task invoked while running make
all: clean .build/dblite.so

# target to build a dynamic extension that can be loaded at runtime
.build/dblite.so: $(shell find . -type f -name '*.go' -o -name '*.c')
	$(call log, $(CYAN), "building $@")
	@go build -buildmode=c-shared -o $@ .
	$(call log, $(GREEN), "built $@")

clean:
	$(call log, $(YELLOW), "nuking .build/")
	@-rm -rf .build/

# ========================================
# target for common golang tasks

# go build tags used by test, vet and more
TAGS = "static,system_libgit2"

update:
	go get -tags=$(TAGS) -u ./...

vet:
	go vet -v -tags=$(TAGS) ./...

build:
	go build -v -tags=$(TAGS) mergestat.go

lint:
	golangci-lint run --build-tags $(TAGS)

lint-ci:
	./bin/golangci-lint run --build-tags $(TAGS) --out-format github-actions

test:
	go test -v -tags=$(TAGS) ./...

test-cover:
	go test -v -tags=$(TAGS) ./... -cover -covermode=count -coverprofile=coverage.out
	go tool cover -html=coverage.out

bench:
	go test -v -tags=$(TAGS) -bench=. -benchmem -run=^nomatch ./...

# ========================================
# some utility methods

# ASCII color codes that can be used with functions that output to stdout
RED		:= 1;31
GREEN	:= 1;32
ORANGE	:= 1;33
YELLOW	:= 1;33
BLUE	:= 1;34
PURPLE	:= 1;35
CYAN	:= 1;36

# log:
#	print out $2 to stdout using $1 as ASCII color codes
define log
	@printf "\033[$(strip $1)m-- %s\033[0m\n" $2
endef
