.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	go generate -x
	embedmd -w README.md

.PHONY: bootstrap
bootstrap:
	GO111MODULE="off" go get -u \
		github.com/campoy/embedmd \
		modernc.org/assets
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin v1.19.1
