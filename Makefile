.PHONY: help clean build fmt lint vet run test style cyclo

SOURCES:=$(shell find . -name '*.go')
DOCFILES:=$(addprefix docs/packages/, $(addsuffix .html, $(basename ${SOURCES})))

default: build

clean: ## Run go clean
	@go clean

build: ## Run go build
	@go build

docs/packages/%.html: %.go
	mkdir -p $(dir $@)
	docgo -outdir $(dir $@) $^

godoc: ${DOCFILES}
