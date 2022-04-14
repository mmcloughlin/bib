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
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin" v1.45.2
