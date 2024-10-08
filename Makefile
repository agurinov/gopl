.PHONY: FORCE
FORCE:

GO               := go
JQ               := jq
GIT              := git
DELVE            := dlv
GOLANGCI_LINT    := golangci-lint
FIELDALIGNMENT   := fieldalignment
GOVULNCHECK      := govulncheck
GO_COVER_TREEMAP := go-cover-treemap

IS_CI := $(firstword $(CI) $(GITLAB_CI) $(GITHUB_ACTIONS) $(CIRCLECI) $(DRONE) $(TEAMCITY_VERSION))

# FHS {{{
# https://en.wikipedia.org/wiki/Filesystem_Hierarchy_Standard
FHS_ROOTDIR := $(CURDIR)/.gura
FHS_BINDIR  := $(FHS_ROOTDIR)/bin
FHS_ETCDIR  := $(FHS_ROOTDIR)/etc
FHS_LIBDIR  := $(FHS_ROOTDIR)/lib

$(FHS_ROOTDIR) $(FHS_BINDIR) $(FHS_ETCDIR) $(FHS_LIBDIR):
	mkdir -p $@
# }}}

# GO public API {{{
GO_TAGS     ?=
GO_PKG      ?= ./...
GO_SYMBOL   ?= .
GO_CMD_ARGS ?=
# }}}

# GO private API {{{
_GO_VERSION_RAW          := $(shell $(GO) env GOVERSION;)
GO_VERSION               := $(subst go,,$(_GO_VERSION_RAW))
_GO_VERSION_SEMVER_PARTS := $(subst ., ,$(GO_VERSION))
_GO_VERSION_SEMVER_MAJOR := $(word 1,$(_GO_VERSION_SEMVER_PARTS))
_GO_VERSION_SEMVER_MINOR := $(word 2,$(_GO_VERSION_SEMVER_PARTS))
_GO_VERSION_SEMVER_PATCH := $(word 3,$(_GO_VERSION_SEMVER_PARTS))
GO_VERSION_MINOR         := $(_GO_VERSION_SEMVER_MAJOR).$(_GO_VERSION_SEMVER_MINOR)

GO_CMD_DIR                 := $(realpath $(CURDIR)/cmd)
GO_CMD_FILES               := $(realpath $(CURDIR)/main.go)
GO_MODULE_FILE             := $(realpath go.mod)
GO_WORKSPACE_FILE          := $(realpath go.work)
GO_MODULE_PATH             := $(shell $(GO) mod edit -json | $(JQ) -Mr '.Module.Path';)
GO_MODULE_GO_VERSION       := $(shell $(GO) mod edit -json | $(JQ) -Mr '.Go';)
GO_MODULE_NAME             := $(basename $(notdir $(GO_MODULE_PATH)))
UNSPECIFIED_GO_MODULE_NAME := github.com/acme/goplay
# }}}

# deps / mod / vendor {{{
go.mod:
	$(eval URL            := $(shell git config --get remote.origin.url;))
	$(eval URL_SCHEMA     := $(firstword $(subst ://, ,$(URL)))://)
	$(eval URL_SCHEMALESS := $(subst $(URL_SCHEMA),,$(URL)))
	$(GO) mod init '$(firstword $(URL_SCHEMALESS) $(UNSPECIFIED_GO_MODULE_NAME))'

go_vendor: FORCE go.mod go_mod_no_cache
ifneq '$(GO_WORKSPACE_FILE)' ''
	$(GO) work vendor
else
	$(GO) mod tidy -v
	$(GO) mod vendor
	$(GO) mod verify
endif
# }}}

# build tags {{{
GREP_TAGS_CMD                     := grep -I -h -R --exclude-dir=.git --exclude-dir=vendor --include=*.go '//go:build' | awk '{for (i=2; i<=NF; i++) print $$i}'

_GO_IGNORED_BUILD_TAGS            := ignore devtools tools neverbuild
_GO_INTEGRATION_BUILD_TAGS        := integration integ
_GO_OS_ARCH_BUILD_TAGS            := unix linux darwin windows 386 amd64 arm arm64 wasm
_GO_SYS_BUILD_TAGS                := cgo gc gccgo

SPECIAL_CHARS                     := ! ( ) & |
PARENTHESIS_OPEN                  := (
PARENTHESIS_CLOSE                 := )

GO_DISCOVERED_BUILD_TAGS_ALL      := $(shell $(GREP_TAGS_CMD);)
GO_DISCOVERED_BUILD_TAGS_ALL      := $(subst !,,$(GO_DISCOVERED_BUILD_TAGS_ALL))
GO_DISCOVERED_BUILD_TAGS_ALL      := $(subst &,,$(GO_DISCOVERED_BUILD_TAGS_ALL))
GO_DISCOVERED_BUILD_TAGS_ALL      := $(subst |,,$(GO_DISCOVERED_BUILD_TAGS_ALL))
GO_DISCOVERED_BUILD_TAGS_ALL      := $(subst $(PARENTHESIS_OPEN),,$(GO_DISCOVERED_BUILD_TAGS_ALL))
GO_DISCOVERED_BUILD_TAGS_ALL      := $(subst $(PARENTHESIS_CLOSE),,$(GO_DISCOVERED_BUILD_TAGS_ALL))
GO_DISCOVERED_BUILD_TAGS_ALL      := $(filter-out $(_GO_OS_ARCH_BUILD_TAGS),$(GO_DISCOVERED_BUILD_TAGS_ALL))
GO_DISCOVERED_BUILD_TAGS_ALL      := $(filter-out $(_GO_SYS_BUILD_TAGS),$(GO_DISCOVERED_BUILD_TAGS_ALL))
GO_DISCOVERED_BUILD_TAGS_ALL      := $(sort $(GO_DISCOVERED_BUILD_TAGS_ALL))
GO_DISCOVERED_BUILD_TAGS_FILTERED := $(filter-out $(_GO_IGNORED_BUILD_TAGS),$(GO_DISCOVERED_BUILD_TAGS_ALL))

ifndef GO_TAGS
go_test:                     GO_TAGS := $(filter-out $(_GO_INTEGRATION_BUILD_TAGS),$(GO_DISCOVERED_BUILD_TAGS_FILTERED))
go_itest dlv_test go_bench:  GO_TAGS := $(GO_DISCOVERED_BUILD_TAGS_FILTERED)
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
TEST_RUN_FLAG     :=

LD_FLAGS    := -s -w
BUILD_FLAGS  = -mod=vendor -tags='$(GO_TAGS)' $(BUILD_RACE_FLAG) -trimpath -ldflags='$(LD_FLAGS)'
COVER_FLAGS := -cover -covermode=atomic
PPROF_FLAGS := -cpuprofile=cpu.pprof -memprofile=mem.pprof -blockprofile=block.pprof -mutexprofile=mutex.pprof -trace=trace.trace
TEST_FLAGS   = $(BUILD_FLAGS) $(COVER_FLAGS) $(TEST_VERBOSE_FLAG) $(TEST_SHORT_FLAG) -run='$(TEST_RUN_FLAG)' -timeout='$(TEST_TIMEOUT_FLAG)' -failfast -count=$(TEST_COUNT_FLAG) $(TEST_USER_FLAGS)
BENCH_FLAGS  = $(BUILD_FLAGS) $(COVER_FLAGS) $(TEST_VERBOSE_FLAG) $(TEST_SHORT_FLAG) -run='$(TEST_RUN_FLAG)' -timeout='$(TEST_TIMEOUT_FLAG)' $(PPROF_FLAGS) -bench=. -benchmem $(TEST_USER_FLAGS)
# }}}

# G_DEBUG / caches {{{
# https://github.com/golang/go/wiki/CoreDumpDebugging
go_mod_no_cache        : FORCE
go_build_no_cache      : FORCE
go_test_no_cache       : FORCE
golangci_lint_no_cache : FORCE

ifdef G_DEBUG
# TEST_VERBOSE_FLAG := -v
TEST_COUNT_FLAG   := 1
TEST_TIMEOUT_FLAG := 5m

go_mod_no_cache        : FORCE ; $(GO) clean -modcache
go_build_no_cache      : FORCE ; $(GO) clean -cache
go_test_no_cache       : FORCE ; $(GO) clean -testcache -fuzzcache
golangci_lint_no_cache : FORCE ; $(GOLANGCI_LINT) cache clean
go_build               : BUILD_FLAGS += -gcflags='-m=2'
go_shared              : BUILD_FLAGS += -gcflags='-m=2'
go_run                 : BUILD_FLAGS += -gcflags='-m=2'
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
	test -r '$(GO_CMD_FILE)'
	$(GO) run $(BUILD_FLAGS) $(GO_CMD_FILE) $(GO_CMD_ARGS)
# }}}

# test / bench / delve {{{
ifndef IS_CI
go_test:  COVER_FLAGS += -coverprofile='$(COVERAGE_OUT_FLAG)'
go_itest: COVER_FLAGS += -coverprofile='$(COVERAGE_OUT_FLAG)'
go_bench: COVER_FLAGS += -coverprofile='$(COVERAGE_OUT_FLAG)'

define COVER_FILES_CMD
	test -r '$(COVERAGE_OUT_FLAG)'
	sed -i '' '/mock/d' '$(COVERAGE_OUT_FLAG)'
	sed -i '' '/gen.go/d' '$(COVERAGE_OUT_FLAG)'
	$(GO) tool cover -html='$(COVERAGE_OUT_FLAG)' -o '$(COVERAGE_FILE_BASENAME).html'
	$(GO) tool cover -func='$(COVERAGE_OUT_FLAG)' -o '$(COVERAGE_FILE_BASENAME).func'
	$(GO_COVER_TREEMAP) -coverprofile '$(COVERAGE_OUT_FLAG)' -statements=false -h 1080 -w 1080 > '$(COVERAGE_FILE_BASENAME).svg'
	test -r '$(COVERAGE_FILE_BASENAME).html'
	test -r '$(COVERAGE_FILE_BASENAME).func'
	test -r '$(COVERAGE_FILE_BASENAME).svg'
endef
endif

go_test: COVERAGE_FILE_BASENAME := coverage_test

go_itest: TEST_SHORT_FLAG        :=
go_itest: TEST_COUNT_FLAG        := 1
go_itest: TEST_TIMEOUT_FLAG      := 5m
go_itest: TEST_RUN_FLAG          := ^Test
go_itest: COVERAGE_FILE_BASENAME := coverage_test

go_test go_itest: FORCE go_test_no_cache
	$(GO) test $(TEST_FLAGS) $(GO_PKG)
	$(COVER_FILES_CMD)

# https://pkg.go.dev/golang.org/x/perf/cmd/benchstat
go_bench: COVERAGE_FILE_BASENAME := coverage_bench
go_bench: TEST_RUN_FLAG          := NONE
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
$(GOLANGCI_LINT) run --config='$(GOLANGCI_LINT_CONF)' --go='$(GO_VERSION_MINOR)' --out-format=colored-line-number --modules-download-mode=vendor --build-tags='$(GO_TAGS)' --tests

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
	$(GOVULNCHECK) -test -tags='$(GO_TAGS)' $(GO_PKG)
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
lint     : FORCE go_vet golangci-lint
fmt      : FORCE go_fmt
generate : FORCE go_generate
vendor   : FORCE go_vendor
