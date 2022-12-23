ROOT_DIR := $(CURDIR)
BUILD_DIR := $(ROOT_DIR)/_output
BIN_DIR := $(BUILD_DIR)/bin
REPO_PREFIX := yunion.io/x/cloudmux
VERSION_PKG := yunion.io/x/pkg/util/version

output_dir:
	@mkdir -p $(BUILD_DIR)

bin_dir: output_dir
	@mkdir -p $(BUILD_DIR)/bin

prepare_dir: bin_dir

GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git branch -r --contains | head -1 | sed -E -e "s%(HEAD ->|origin|upstream)/?%%g" | xargs)
GIT_VERSION := $(shell git describe --always --tags --abbrev=14 $(GIT_COMMIT)^{commit})
GIT_TREE_STATE := $(shell s=`git status --porcelain 2>/dev/null`; if [ -z "$$s" ]; then echo "clean"; else echo "dirty"; fi)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := "-w \
	-X $(VERSION_PKG).gitVersion=$(GIT_VERSION) \
	-X $(VERSION_PKG).gitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PKG).gitBranch=$(GIT_BRANCH) \
	-X $(VERSION_PKG).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PKG).gitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PKG).gitMajor=0 \
	-X $(VERSION_PKG).gitMinor=0"

ifneq ($(DLV),)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
	LDFLAGS = ""
endif
# GO_BUILD_FLAGS+=-mod vendor -ldflags $(LDFLAGS)
GO_BUILD_FLAGS+=-ldflags $(LDFLAGS)
GO_BUILD := go build $(GO_BUILD_FLAGS)

cmd/%: prepare_dir
	$(GO_BUILD) -o $(BIN_DIR)/$(shell basename $@) $(REPO_PREFIX)/$@

test:
	go test -v $(GO_BUILD_FLAGS) ./...

fmt:
	goimports -w -local "yunion.io/x/:yunion.io/x/onecloud:yunion.io/x/cloudmux" pkg cmd

GOPROXY ?= direct

mod:
	GOPROXY=$(GOPROXY) go get -d $(patsubst %,%@master,$(shell GO111MODULE=on go mod edit -print  | sed -n -e 's|.*\(yunion.io/x/[a-z].*\) v.*|\1|p'))
	go mod tidy
