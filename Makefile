BUILD_VERSION   ?= $(shell cat version.txt || echo "0.1")
BUILD_DATE      := $(shell date "+%F %T")
GIT_COMMIT     := $(shell git rev-parse HEAD || echo "0.0.0")

help: ## this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## fmt
	gofmt -w .
	goimports -w .
	@echo gofmt -l
	@OUTPUT=`gofmt -l . 2>&1`; \
	if [ "$$OUTPUT" ]; then \
		echo "gofmt must be run on the following files:"; \
        echo "$$OUTPUT"; \
        exit 1; \
    fi

golint: ## lint

	@echo golangci-lint run --skip-files \".*test.go\" -v ./...
	@OUTPUT=`command -v golangci-lint >/dev/null 2>&1 && golangci-lint run --skip-files ".*test.go"  -v ./... 2>&1`; \
	if [ "$$OUTPUT" ]; then \
		echo "golint errors:"; \
		echo "$$OUTPUT"; \
	fi

default: fmt golint ## fmt code

build: clean generate ## 构建二进制
	@echo "build bin ${BUILD_VERSION} ${BUILD_DATE} ${GIT_COMMIT}"
	# go install github.com/mitchellh/gox@latest
	@gox -osarch="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64" \
        -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" \
    	-ldflags   "-w -s \
    				-X 'github.com/ysicing/ergo/version.Version=${BUILD_VERSION}' \
                    -X 'github.com/ysicing/ergo/version.BuildDate=${BUILD_DATE}' \
                        -X 'github.com/ysicing/ergo/version.GitCommitHash=${GIT_COMMIT}' \
										    -X 'k8s.io/client-go/pkg/version.gitVersion=${BUILD_VERSION}' \
    										-X 'k8s.io/client-go/pkg/version.gitCommit=${GIT_COMMIT}' \
    										-X 'k8s.io/client-go/pkg/version.gitTreeState=dirty' \
    										-X 'k8s.io/client-go/pkg/version.buildDate=${BUILD_DATE}' \
												-X 'k8s.io/client-go/pkg/version.gitMajor=1' \
												-X 'k8s.io/client-go/pkg/version.gitMinor=23' \
    										-X 'k8s.io/component-base/version.gitVersion=${BUILD_VERSION}' \
    										-X 'k8s.io/component-base/version.gitCommit=${GIT_COMMIT}' \
                        -X 'k8s.io/component-base/version.gitTreeState=dirty' \
												-X 'k8s.io/component-base/version.gitMajor=1' \
												-X 'k8s.io/component-base/version.gitMinor=23' \
    										-X 'k8s.io/component-base/version.buildDate=${BUILD_DATE}'"

upx: ## upx binary
	@echo "upx binray"
	@upx -9 dist/*

generate: ## generate
	go generate ./...

docker: ## 构建镜像
	@echo "build docker images ${BUILD_VERSION}"
	@docker build -t ysicing/ergo .
	@docker build -t ysicing/ergo:${BUILD_VERSION} .

dpush: docker
	@docker push ysicing/ergo
    @docker push ysicing/ergo:${BUILD_VERSION}

release:  ## github release
	ghr -u ysicing -t $(GITHUB_RELEASE_TOKEN) -b "release ${BUILD_VERSION}" -n "release ${BUILD_VERSION}" -soft --debug ${BUILD_VERSION} dist

pre-release:  ## github pre-release
	ghr -u ysicing -t $(GITHUB_RELEASE_TOKEN) -b "release ${BUILD_VERSION}" -n "release ${BUILD_VERSION}" -replace -recreate -prerelease --debug ${BUILD_VERSION} dist

clean: ## clean
	rm -rf dist

install: clean ## install
	go install \
		-ldflags   "-w -s \
						-X 'github.com/ysicing/ergo/version.Version=${BUILD_VERSION}' \
                        -X 'github.com/ysicing/ergo/version.BuildDate=${BUILD_DATE}' \
                        -X 'github.com/ysicing/ergo/version.GitCommitHash=${GIT_COMMIT}' \
										    -X 'k8s.io/client-go/pkg/version.gitVersion=${BUILD_VERSION}' \
    										-X 'k8s.io/client-go/pkg/version.gitCommit=${GIT_COMMIT}' \
    										-X 'k8s.io/client-go/pkg/version.gitTreeState=dirty' \
    										-X 'k8s.io/client-go/pkg/version.buildDate=${BUILD_DATE}' \
    										-X 'k8s.io/component-base/version.gitVersion=${BUILD_VERSION}' \
    										-X 'k8s.io/component-base/version.gitCommit=${GIT_COMMIT}' \
                        -X 'k8s.io/component-base/version.gitTreeState=dirty' \
    										-X 'k8s.io/component-base/version.buildDate=${BUILD_DATE}'"

deb: build ## build deb
	./deb.sh

doc: ## gen docs
	go run ./hack/gendoc/doc.go
	cp -a docs/ergo.md docs/index.md
	cp -a README.md docs/readme.md

# cleanvm: ## clem lima vm
# 	limactl ls | grep debian && (limactl stop debian || echo "skip stop") &&limactl rm debian || echo "not found"

# vm: cleanvm ## start lima vm
# 	limactl start common/debian.yml

# shell: ## shell debian
# 	limactl shell debian

# local-test: build ## 本地测试
# 	limactl cp ./dist/ergo_linux_amd64 lima-ergo:/home/ysicing.linux/.ergo/bin/ergo-ergo

local: build ## local test
	upx -1 ./dist/ergo_linux_amd64
	scp ./dist/ergo_linux_amd64 192.168.0.107:/root/

.PHONY : build release clean install

.EXPORT_ALL_VARIABLES:

GO111MODULE = on
GOPROXY = https://goproxy.cn
GOSUMDB = sum.golang.google.cn
