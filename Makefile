TAG=$(shell git describe --tags --always)

all: | test vet lint

test:
	go test ./...

vet:
	go vet ./...

lint:
	cd /tmp && go get golang.org/x/lint/golint
	golint ./...

announce:
	GOPROXY=proxy.golang.org go list -m github.com/and3rson/telemux@${TAG}
	http https://sum.golang.org/lookup/github.com/and3rson/telemux@${TAG}
	http https://proxy.golang.org/github.com/and3rson/telemux/@v/${TAG}.info
	cd /tmp && mkdir -p .go && chmod -R 777 .go && rm -rf .go && GOPATH=/tmp/.go GOPROXY=https://proxy.golang.org GO111MODULE=on go get github.com/and3rson/telemux@${TAG}

changelog:
	./mkchangelog.sh > ./CHANGELOG.md
