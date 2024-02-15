# Description: Makefile for building the wasm file

WASM_EXEC_TINY := $(shell tinygo env TINYGOROOT)/targets/wasm_exec.js

build:
	cp "$(WASM_EXEC_TINY)" .
	GOOS=darwin tinygo build -o main.wasm

install:
	brew tap tinygo-org/tools
	brew install tinygo

