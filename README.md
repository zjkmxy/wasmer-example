Requirements
============

One of the following compilers is enough.

Install WASM-SDK and WASMER
-------------

Wasienv seems to be abandoned. Use the following command in a setup folder to download MacOS builds.

```bash
export WASI_VERSION=20
export WASI_VERSION_FULL=${WASI_VERSION}.0
wget https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-${WASI_VERSION}/wasi-sdk-${WASI_VERSION_FULL}-macos.tar.gz
tar xvf wasi-sdk-${WASI_VERSION_FULL}-macos.tar.gz
```

Then, in `bashrc`, export `WASISDK_DIR` to be the `bin` folder.


Install wasmer via Homebrew `brew install wabt wasmer`.

To check:
```bash
${WASISDK_DIR}/clang++ -v
# clang version 16.0.0 (https://github.com/llvm/llvm-project 08d094a0e457360ad8b94b017d2dc277e697ca76)
# Target: wasm32-unknown-wasi
# Thread model: posix

wasmer --version
# wasmer 3.3.0
```

Now we can use `wasmc++` and `wasic++` to compile C++ sources into WASM/WASI.

Install [Emscripten](https://emscripten.org/docs/getting_started/downloads.html)
-------------

Install with `brew install emscripten`

To check:
```bash
emcc --version
# emcc (Emscripten gcc/clang-like replacement + linker emulating GNU ld) 3.1.39-git
# Copyright (C) 2014 the Emscripten authors (see AUTHORS.txt)
# This is free and open source software under the MIT license.
# There is NO warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
```

Differences between compilers
-----------------------------

- `wasic++` will include the standard library.
- `wasmc++` will not include any extra library. The generated wasm will be very small in size.
  - Note: `wasic++` will also include syscalls, but by default Wasmer does not expose them.
- Emscripten is more "on the browser's side".

Compilation
===========

Execute `make` in `csrc` folder to compile C++ code into WASM.
Then you can execute the main program by `go run .`.
