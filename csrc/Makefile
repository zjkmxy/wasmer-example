WASI_COMPILER=${WASISDK_DIR}/clang++
WASM_FLAG=-Wl,--allow-undefined -Os
WASI_FLAG=$(WASM_FLAG) -Wl,--no-entry

all: sample.wasm

clean:
	rm *.wasm

sample.wasm: sample.cpp
	$(WASI_COMPILER) $(WASI_FLAG) $< -o $@

sample2.wasm: sample2.wat
	wat2wasm $< -o $@

sample_emcc.wasm: sample_emcc.cpp
	emcc $< -Os -s STANDALONE_WASM -s ERROR_ON_UNDEFINED_SYMBOLS=0 -o $@
