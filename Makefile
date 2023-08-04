.DEFAULT_GOAL := all

.PHONY: all
all: tidy gen add-copyright format lint cover build


fmt:
	@command -v gofumpt || (WORK=$(shell pwd) && cd /tmp && GO111MODULE=on go install mvdan.cc/gofumpt@latest && cd $(WORK))
	@gofumpt -w  -d .
	@command -v gci || go install github.com/daixiang0/gci@latest
	@gci write -s standard -s default -s 'Prefix(github.com/kart-io/)' --skip-generated .

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint:
	golangci-lint version
	golangci-lint run -v --color always --out-format colored-line-number