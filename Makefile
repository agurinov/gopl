.PHONY: FORCE
FORCE:

include gitlabci.mk

GO             := go
GIT            := git
DELVE          := dlv
GOLANGCI_LINT  := golangci-lint
FIELDALIGNMENT := fieldalignment
GOVULNCHECK    := govulncheck

IS_CI := $(firstword $(CI) $(GITLAB_CI) $(GITHUB_ACTIONS) $(CIRCLECI) $(DRONE))

# FHS {{{
# https://en.wikipedia.org/wiki/Filesystem_Hierarchy_Standard
FHS_ROOTDIR := $(CURDIR)/.gura
FHS_BINDIR  := $(FHS_ROOTDIR)/bin
FHS_ETCDIR  := $(FHS_ROOTDIR)/etc
FHS_LIBDIR  := $(FHS_ROOTDIR)/lib

$(FHS_ROOTDIR) $(FHS_BINDIR) $(FHS_ETCDIR) $(FHS_LIBDIR):
	mkdir -p $@
# }}}

# GO private API {{{
GO_VERSION_RAW          := $(shell $(GO) env GOVERSION;)
GO_VERSION              := $(subst go,,$(GO_VERSION_RAW))
GO_VERSION_SEMVER_PARTS := $(subst ., ,$(GO_VERSION))
GO_VERSION_SEMVER_MAJOR := $(word 1,$(GO_VERSION_SEMVER_PARTS))
GO_VERSION_SEMVER_MINOR := $(word 2,$(GO_VERSION_SEMVER_PARTS))
GO_VERSION_SEMVER_PATCH := $(word 3,$(GO_VERSION_SEMVER_PARTS))
GO_VERSION_MINOR        := $(GO_VERSION_SEMVER_MAJOR).$(GO_VERSION_SEMVER_MINOR)

GO_CMD_DIR              := $(realpath $(CURDIR)/cmd)
GO_CMD_FILES            := $(realpath $(CURDIR)/main.go)
GO_MODULE_FILE          := $(realpath go.mod)
GO_MODULE_PATH          := $(shell $(GO) mod edit -json | jq -Mr '.Module.Path';)
GO_MODULE_GO_VERSION    := $(shell $(GO) mod edit -json | jq -Mr '.Go';)
GO_MODULE_NAME          := $(basename $(notdir $(GO_MODULE_PATH)))
# }}}

# GO public API {{{
GO_TAGS     ?=
GO_PKG      ?= ./...
GO_CMD_ARGS ?=
# }}}

# deps / mod / vendor {{{
go.mod:
	$(GO) mod init '$(URL_SCHEMALESS)'

go_vendor: FORCE go.mod go_mod_no_cache
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify
# }}}

# build tags {{{
GREP_TAGS_CMD                     := grep -I -h -R --exclude-dir=.git --exclude-dir=vendor --include=*.go '//go:build' | awk '{print $$2}'

GO_IGNORED_BUILD_TAGS             := ignore devtools neverbuild
GO_DISCOVERED_BUILD_TAGS_ALL      := $(sort $(shell $(GREP_TAGS_CMD);))
GO_DISCOVERED_BUILD_TAGS_FILTERED := $(filter-out $(GO_IGNORED_BUILD_TAGS),$(GO_DISCOVERED_BUILD_TAGS_ALL))

ifndef GO_TAGS
go_test:                     GO_TAGS := test_unit
dlv_test go_bench:           GO_TAGS := $(GO_DISCOVERED_BUILD_TAGS_FILTERED)
golangci-lint go_vet go_sec: GO_TAGS := $(GO_DISCOVERED_BUILD_TAGS_FILTERED)
go_generate:                 GO_TAGS := $(GO_DISCOVERED_BUILD_TAGS_ALL)
endif
# }}}

# cmd {{{
ifndef GO_CMD_FILES
ifdef GO_CMD_DIR
FIND_CMD     := find $(GO_CMD_DIR) -type f -name
GO_CMD_FILES += $(realpath \
	$(shell $(FIND_CMD) main.go;) \
	$(shell $(FIND_CMD) cmd.go;) \
	$(shell $(FIND_CMD) $(GO_MODULE_NAME).go;) \
)
endif
endif

GO_CMD_FILES := $(sort $(GO_CMD_FILES))
# }}}

# CLI {{{
BUILD_RACE_FLAG   := -race
COVERAGE_OUT_FLAG := coverage.out
TEST_VERBOSE_FLAG :=
TEST_SHORT_FLAG   := -short
TEST_COUNT_FLAG   := $(if $(findstring ./...,$(GO_PKG)),1,20)
TEST_TIMEOUT_FLAG := 5s

BUILD_FLAGS  = -mod=vendor -tags='$(GO_TAGS)' $(BUILD_RACE_FLAG) -trimpath
COVER_FLAGS := -cover -covermode=atomic
PPROF_FLAGS := -cpuprofile=cpu.pprof -memprofile=mem.pprof -blockprofile=block.pprof -mutexprofile=mutex.pprof -trace=trace.trace
TEST_FLAGS   = $(BUILD_FLAGS) $(COVER_FLAGS) $(TEST_VERBOSE_FLAG) $(TEST_SHORT_FLAG) -failfast -timeout='$(TEST_TIMEOUT_FLAG)' -count=$(TEST_COUNT_FLAG) $(TEST_USER_FLAGS)
BENCH_FLAGS  = $(BUILD_FLAGS) $(COVER_FLAGS) $(TEST_VERBOSE_FLAG) $(TEST_SHORT_FLAG) $(PPROF_FLAGS) -bench=. -benchmem -run=NONE
# }}}

# G_DEBUG / caches {{{
# https://github.com/golang/go/wiki/CoreDumpDebugging
go_mod_no_cache        : FORCE
go_build_no_cache      : FORCE
go_test_no_cache       : FORCE
golangci_lint_no_cache : FORCE

ifdef G_DEBUG
TEST_VERBOSE_FLAG := -v
TEST_SHORT_FLAG   :=
TEST_COUNT_FLAG   := 1
TEST_TIMEOUT_FLAG := 5m

go_mod_no_cache        : FORCE ; $(GO) clean -modcache
go_build_no_cache      : FORCE ; $(GO) clean -cache
go_test_no_cache       : FORCE ; $(GO) clean -testcache -fuzzcache
golangci_lint_no_cache : FORCE ; $(GOLANGCI_LINT) cache clean

go_build  : BUILD_FLAGS += -gcflags='-m=2'
go_shared : BUILD_FLAGS += -gcflags='-m=2'
go_run    : BUILD_FLAGS += -gcflags='-m=2'
endif
# }}}

# build / run {{{
$(GO_CMD_FILES): FORCE $(FHS_BINDIR)
	$(eval GO_CMD_PKG  := $(realpath $(dir $@)))
	$(eval GO_CMD_NAME := $(basename $(notdir $(GO_CMD_PKG))))
	$(GO) build $(BUILD_FLAGS) -o $(FHS_BINDIR)/$(GO_CMD_NAME) $(GO_CMD_PKG)

go_build: FORCE vendor go_build_no_cache $(GO_CMD_FILES)

go_shared: GO_CMD_FILE := shared.go
go_shared: BUILD_FLAGS += -buildmode=c-shared
go_shared: FORCE vendor go_build_no_cache $(FHS_LIBDIR) $(GO_CMD_FILE)
	$(eval GO_CMD_PKG  := $(realpath $(dir $@)))
	$(eval GO_CMD_NAME := $(basename $(notdir $(GO_CMD_PKG))))
	$(GO) build $(BUILD_FLAGS) -o $(FHS_LIBDIR)/$(GO_CMD_NAME).so $(GO_CMD_FILE)

go_run: GO_CMD_FILE := $(firstword $(GO_CMD_FILES))
go_run: FORCE vendor go_build_no_cache
	test -r $(GO_CMD_FILE)
	$(GO) run $(BUILD_FLAGS) $(GO_CMD_FILE) $(GO_CMD_ARGS)
# }}}

# test / bench / delve {{{
ifndef IS_CI
go_test:  COVER_FLAGS += -coverprofile='$(COVERAGE_OUT_FLAG)'
go_bench: COVER_FLAGS += -coverprofile='$(COVERAGE_OUT_FLAG)'

define COVER_FILES_CMD
	test -r '$(COVERAGE_OUT_FLAG)'
	$(GO) tool cover -html='$(COVERAGE_OUT_FLAG)' -o '$(COVERAGE_FILE_BASENAME).html'
	$(GO) tool cover -func='$(COVERAGE_OUT_FLAG)' -o '$(COVERAGE_FILE_BASENAME).func'
	test -r '$(COVERAGE_FILE_BASENAME).html'
	test -r '$(COVERAGE_FILE_BASENAME).func'
endef
endif

go_test: COVERAGE_FILE_BASENAME := coverage_test
go_test: FORCE go_test_no_cache
	$(GO) test $(TEST_FLAGS) $(GO_PKG)
	$(COVER_FILES_CMD)

# https://pkg.go.dev/golang.org/x/perf/cmd/benchstat
go_bench: COVERAGE_FILE_BASENAME := coverage_bench
go_bench: FORCE go_test_no_cache
	$(GO) test $(BENCH_FLAGS) $(GO_PKG)
	$(COVER_FILES_CMD)

dlv_test: BUILD_RACE_FLAG :=
dlv_test: TEST_COUNT_FLAG := 1
dlv_test: COVER_FLAGS     :=
dlv_test: FORCE go_test_no_cache
	$(DELVE) test --build-flags='$(TEST_FLAGS)' $(GO_PKG)
# }}}

# lint {{{

# https://golangci-lint.run/usage/configuration/#config-file
GOLANGCI_LINT_CONFS := $(realpath \
	$(CURDIR)/.golangci.yml \
	$(CURDIR)/.golangci.yaml \
	$(CURDIR)/.golangci.toml \
	$(CURDIR)/.golangci.json \
	$(FHS_ETCDIR)/.golangci.yaml \
)

define GOLANGCI_LINT_CMD
$(GOLANGCI_LINT) run --go='$(GO_VERSION_MINOR)' --modules-download-mode=vendor --build-tags='$(GO_TAGS)' --config='$(GOLANGCI_LINT_CONF)'

endef

golangci-lint: FORCE golangci_lint_no_cache | $(firstword $(GOLANGCI_LINT_CONFS))
	$(foreach GOLANGCI_LINT_CONF,$|,$(GOLANGCI_LINT_CMD))

go_fmt: FORCE
	$(GO) fmt $(GO_PKG)
	$(GO) mod edit -fmt
	$(FIELDALIGNMENT) -fix $(GO_PKG)

go_vet: FORCE
	$(GO) vet $(BUILD_FLAGS) -c=5 $(GO_PKG)

go_sec: FORCE
	-$(GOVULNCHECK) -test -tags='$(GO_TAGS)' $(GO_PKG)
# }}}

# codegen {{{
go_generate: FORCE
	$(GO) generate -x $(BUILD_FLAGS) $(GO_PKG)
# }}}

# git {{{
git_clean: FORCE
	$(GIT) clean -ffdx
# }}}

# environment info {{{
GO_DOT_VARIABLES   := $(sort $(filter GO_%,$(.VARIABLES)))
$(GO_DOT_VARIABLES): FORCE ; @echo $@='$($@)'
go_info: FORCE $(GO_DOT_VARIABLES)
	@$(GO) env
# }}}

.DEFAULT_GOAL := go_info

clean    : FORCE go_mod_no_cache go_build_no_cache go_test_no_cache golangci_lint_no_cache git_clean
lint     : FORCE go_vet go_sec golangci-lint
fmt      : FORCE go_fmt
generate : FORCE go_generate
vendor   : FORCE go_vendor
