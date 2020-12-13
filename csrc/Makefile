WASM_COMPILER=wasmc++
WASI_COMPILER=wasic++
WASM_FLAG=-Wl,--allow-undefined
WASI_FLAG=$(WASM_FLAG) -Wl,--export-all -Wl,--no-entry

all: sample.wasm

sample.wasm: sample.cpp
	$(WASI_COMPILER) $(WASI_FLAG) $< -o $@
