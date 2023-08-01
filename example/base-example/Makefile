PROJECTNAME=\$(shell basename "\$(PWD)")
TCR_HOST_LOCAL := costa92

.PHONY: wire
wire:
	wire ./...

.PHONY: run
run:
	@go run -race main.go wire_gen.go
	@echo "LS_OUTPUT="$(LS_OUTPUT);
	@for item  in $(LS_OUTPUT); \
	do \
	  @echo item
	done

vet:
	go vet

lint:
	golangci-lint run --enable-all

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'