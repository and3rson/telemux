TAG ?= $(shell git tag --points-at HEAD | grep -v gormpersistence)

all: | init test vet lint

init:
	cd /tmp && go get golang.org/x/lint/golint
	cd /tmp && go get golang.org/x/tools/cmd/cover

test:
	go test ./... -cover -coverprofile c.out -test.v
	go tool cover -html=c.out -o cover.html
	make -C gormpersistence test

vet:
	go vet ./...

lint:
	golint ./...

announce:
	GOPROXY=proxy.golang.org go list -m github.com/and3rson/telemux@${TAG}
	http https://sum.golang.org/lookup/github.com/and3rson/telemux@${TAG}
	http https://proxy.golang.org/github.com/and3rson/telemux/@v/${TAG}.info
	cd /tmp && mkdir -p .go && chmod -R 777 .go && rm -rf .go && GOPATH=/tmp/.go GOPROXY=https://proxy.golang.org GO111MODULE=on go get github.com/and3rson/telemux@${TAG}
	make -C gormpersistence announce

changelog:
	./mkchangelog.sh > ./CHANGELOG.md
