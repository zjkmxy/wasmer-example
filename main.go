package main

import (
	"fmt"
	"os"

	"github.com/wasmerio/wasmer-go/wasmer"
)

const WASM_PROG = "./csrc/sample.wasm"

type WasmEnv struct {
	inst    *wasmer.Instance
	measure map[string]interface{}
	memory  *wasmer.Memory

	allocate             wasmer.NativeFunction
	deallocate           wasmer.NativeFunction
	afterReceiveInterest wasmer.NativeFunction
}

// Given we have multiple instances in YaNFD, we need to figure out which one
// A global variable is not a good idea

func ForwardInterest(env interface{}, args []wasmer.Value) ([]wasmer.Value, error) {
	faceID := args[0].I32()
	fmt.Printf("Forwarded to: %d\n", faceID)
	return []wasmer.Value{}, nil
}

func getWasmStr(memData []byte, ptr int32) string {
	ed := int(ptr)
	for i := int(ptr); i < len(memData); i++ {
		if memData[i] == 0 {
			ed = i
			break
		}
	}
	return string(memData[ptr:ed])
}

func GetMeasurementInt(env interface{}, args []wasmer.Value) ([]wasmer.Value, error) {
	name := args[0].I32()
	measure := env.(*WasmEnv).measure
	memory := env.(*WasmEnv).memory
	nameStr := getWasmStr(memory.Data(), name)

	// Will panic if it is not an int
	result := measure[nameStr].(int32)

	return []wasmer.Value{wasmer.NewI32(result)}, nil
}

func SetMeasurementInt(env interface{}, args []wasmer.Value) ([]wasmer.Value, error) {
	name := args[0].I32()
	val := args[1].I32()
	measure := env.(*WasmEnv).measure
	memory := env.(*WasmEnv).memory
	nameStr := getWasmStr(memory.Data(), name)

	measure[nameStr] = val

	return []wasmer.Value{}, nil
}

func AddToMeasurementInt(env interface{}, args []wasmer.Value) ([]wasmer.Value, error) {
	name := args[0].I32()
	val := args[1].I32()
	measure := env.(*WasmEnv).measure
	memory := env.(*WasmEnv).memory
	nameStr := getWasmStr(memory.Data(), name)

	cur := measure[nameStr].(int)
	measure[nameStr] = cur + int(val)

	return []wasmer.Value{}, nil
}

func GetMeasurementBytes(env interface{}, args []wasmer.Value) ([]wasmer.Value, error) {
	name := args[0].I32()
	buf := args[1].I32()
	buflen := args[2].I32()
	measure := env.(*WasmEnv).measure
	memory := env.(*WasmEnv).memory
	nameStr := getWasmStr(memory.Data(), name)

	val := measure[nameStr].([]byte)
	ret := int32(copy(memory.Data()[buf:buf+buflen], val))

	return []wasmer.Value{wasmer.NewI32(ret)}, nil
}

func SetMeasurementBytes(env interface{}, args []wasmer.Value) ([]wasmer.Value, error) {
	name := args[0].I32()
	value := args[1].I32()
	size := args[2].I32()
	measure := env.(*WasmEnv).measure
	memory := env.(*WasmEnv).memory
	nameStr := getWasmStr(memory.Data(), name)

	data := make([]byte, size)
	copy(data, memory.Data()[value:value+size])
	measure[nameStr] = data

	return []wasmer.Value{}, nil
}

func callAfterReceiveInterest(env *WasmEnv, name string) {
	// Allocate the string in WASM memory
	// len(name) is the length in bytes, not in characters.
	allocateResult, _ := env.allocate(int32(len(name)))
	inputPointer := allocateResult.(int32)
	defer env.deallocate(inputPointer)

	ptr := int(inputPointer)
	copy(env.memory.Data()[ptr:ptr+len(name)], name)
	env.memory.Data()[ptr+len(name)] = 0

	_, err := env.afterReceiveInterest(inputPointer, int32(len(name)))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// Create an Engine
	engine := wasmer.NewEngine()

	// Create a Store
	store := wasmer.NewStore(engine)

	fmt.Println("Compiling module...")
	wasmBytes, _ := os.ReadFile(WASM_PROG)
	module, err := wasmer.NewModule(store, wasmBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	measure := make(map[string]interface{})
	measure["counter"] = 0
	measure["hash"] = []byte{0, 0, 0, 0}
	env := &WasmEnv{
		measure: measure,
	}

	importObj := wasmer.NewImportObject()
	importObj.Register(
		"env",
		map[string]wasmer.IntoExtern{
			"GetMeasurementInt": wasmer.NewFunctionWithEnvironment(
				store,
				wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32),
					wasmer.NewValueTypes(wasmer.I32)),
				env,
				GetMeasurementInt,
			),
			"GetMeasurementBytes": wasmer.NewFunctionWithEnvironment(
				store,
				wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32),
					wasmer.NewValueTypes(wasmer.I32)),
				env,
				GetMeasurementBytes,
			),
			"ForwardInterest": wasmer.NewFunctionWithEnvironment(
				store,
				wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32),
					wasmer.NewValueTypes()),
				env,
				ForwardInterest,
			),
			"SetMeasurementInt": wasmer.NewFunctionWithEnvironment(
				store,
				wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32, wasmer.I32),
					wasmer.NewValueTypes()),
				env,
				SetMeasurementInt,
			),
			"AddToMeasurementInt": wasmer.NewFunctionWithEnvironment(
				store,
				wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32, wasmer.I32),
					wasmer.NewValueTypes()),
				env,
				AddToMeasurementInt,
			),
			"SetMeasurementBytes": wasmer.NewFunctionWithEnvironment(
				store,
				wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32),
					wasmer.NewValueTypes()),
				env,
				SetMeasurementBytes,
			),
		},
	)
	// Used by WASI
	importObj.Register("wasi_snapshot_preview1",
		map[string]wasmer.IntoExtern{
			"proc_exit": wasmer.NewFunction(
				store,
				wasmer.NewFunctionType(wasmer.NewValueTypes(wasmer.I32), wasmer.NewValueTypes()),
				func(args []wasmer.Value) ([]wasmer.Value, error) {
					fmt.Printf("WASM quited with %d\n", args[0].I32())
					return []wasmer.Value{}, nil
				},
			),
		})

	// Instantiates the WebAssembly module.
	instance, err := wasmer.NewInstance(module, importObj)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer instance.Close()
	env.inst = instance

	env.memory, err = instance.Exports.GetMemory("memory")
	if err != nil {
		fmt.Println(err)
		return
	}
	env.allocate, err = instance.Exports.GetFunction("Allocate")
	if err != nil {
		fmt.Println(err)
		return
	}
	env.deallocate, err = instance.Exports.GetFunction("Deallocate")
	if err != nil {
		fmt.Println(err)
		return
	}
	env.afterReceiveInterest, err = instance.Exports.GetFunction("AfterReceiveInterest")
	if err != nil {
		fmt.Println(err)
		return
	}

	callAfterReceiveInterest(env, "hello")
	callAfterReceiveInterest(env, "world")

	fmt.Printf("counter = %v\n", measure["counter"])
	fmt.Printf("hash = %x\n", measure["hash"])
}
