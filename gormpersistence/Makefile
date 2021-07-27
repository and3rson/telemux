TAG ?= $(shell git tag --points-at HEAD | grep gormpersistence | sed 's/v2\/gormpersistence\///g' | sed 's/gormpersistence\///g')

test:
	go test ./...

announce:
	GOPROXY=proxy.golang.org go list -m github.com/and3rson/telemux/gormpersistence/v2@${TAG}
	http https://sum.golang.org/lookup/github.com/and3rson/telemux/gormpersistence/v2@${TAG}
	http https://proxy.golang.org/github.com/and3rson/telemux/gormpersistence/v2/@v/${TAG}.info
	cd /tmp && mkdir -p .go && chmod -R 777 .go && rm -rf .go && GOPATH=/tmp/.go GOPROXY=https://proxy.golang.org GO111MODULE=on go get github.com/and3rson/telemux/gormpersistence/v2@${TAG}
