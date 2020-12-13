Requirements
============

One of the following compilers is enough.

Install [Wasienv](https://github.com/wasienv/wasienv)
-------------

The official install script does not work sometimes.
If so, remove `--install-option="--install-scripts=$INSTALL_DIRECTORY/bin"` from the script and run it again.
It will install `wasienv` Python package, print an error and stop.
After that, use the following command to copy the compiler scripts manually:

```bash
cp `python3 -m site --user-base`/bin/was* ~/.wasienv/bin
```

Then, run the install script one more time to make sure it installs a compiler and a runtime for you.

To check:
```bash
wasic++ --version
# wasienv (wasienv gcc/clang-like replacement)

wasmer --version
# wasmer 1.0.0-beta1
```

Now we can use `wasmc++` and `wasic++` to compile C++ sources into WASM/WASI.

Install [Emscripten](https://emscripten.org/docs/getting_started/downloads.html)
-------------

I'm not using this one. Homebrew will work I guess.
This toolchain is more mature than WASI.

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
