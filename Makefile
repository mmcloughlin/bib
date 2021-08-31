.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	go generate -x
	embedmd -w README.md

.PHONY: bootstrap
bootstrap:
	go install github.com/campoy/embedmd@v1.0.0
	go install modernc.org/assets@v1.0.0
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin v1.19.1
