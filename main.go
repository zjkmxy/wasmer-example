package main

/*
#include <stdlib.h>

extern int32_t GetMeasurementInt(void *context, int name);
extern int32_t GetMeasurementBytes(void *context, int name, int buf, int buflen);
extern void SetMeasurementInt(void *context, int name, int value);
extern void AddToMeasurementInt(void *context, int name, int value);
extern void SetMeasurementBytes(void *context, int name, int value, int size);
extern void ForwardInterest(void *context, int faceID);
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/wasmerio/wasmer-go/wasmer"
)

var measure map[string]interface{}
var instance wasmer.Instance

// Given we have multiple instances in YaNFD, we need to figure out which one
// A global variable is not a good idea

//export ForwardInterest
func ForwardInterest(context unsafe.Pointer, faceID int32) {
	// I don't quite understand what is "context". It only occurs in some examples.
	fmt.Printf("Forwarded to: %d\n", faceID)
}

func getWasmStr(ptr int32) string {
	return C.GoString((*C.char)(unsafe.Pointer(&instance.Memory.Data()[int(ptr)])))
}

//export GetMeasurementInt
func GetMeasurementInt(context unsafe.Pointer, name int32) int32 {
	nameStr := getWasmStr(name)
	// Will panic if it is not an int
	return measure[nameStr].(int32)
}

//export SetMeasurementInt
func SetMeasurementInt(context unsafe.Pointer, name int32, val int32) {
	nameStr := getWasmStr(name)
	measure[nameStr] = val
}

//export AddToMeasurementInt
func AddToMeasurementInt(context unsafe.Pointer, name int32, val int32) {
	nameStr := getWasmStr(name)
	cur := measure[nameStr].(int)
	measure[nameStr] = cur + int(val)
}

func pointerToSlice(ptr unsafe.Pointer, length int) (ret []byte) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	header.Data = uintptr(ptr)
	header.Cap = length
	header.Len = length
	return
}

//export GetMeasurementBytes
func GetMeasurementBytes(context unsafe.Pointer, name int32, buf int32, buflen int32) int32 {
	nameStr := getWasmStr(name)
	val := measure[nameStr].([]byte)
	bufPtr := unsafe.Pointer(&instance.Memory.Data()[int(buf)])
	return int32(copy(pointerToSlice(bufPtr, int(buflen)), val))
}

//export SetMeasurementBytes
func SetMeasurementBytes(context unsafe.Pointer, name int32, value int32, size int32) {
	nameStr := getWasmStr(name)
	bufPtr := unsafe.Pointer(&instance.Memory.Data()[int(value)])
	// C.GoBytes copy the content of bufPtr
	data := C.GoBytes(bufPtr, C.int(size))
	measure[nameStr] = data
}

func callAfterReceiveInterest(name string) {
	// Allocate the string in WASM memory
	// len(name) is the length in bytes, not in characters.
	allocateResult, _ := instance.Exports["Allocate"](len(name))
	inputPointer := allocateResult.ToI32()
	defer instance.Exports["Deallocate"](inputPointer)

	strPtr := unsafe.Pointer(&instance.Memory.Data()[int(inputPointer)])
	copy(pointerToSlice(strPtr, len(name)), name)

	_, err := instance.Exports["AfterReceiveInterest"](inputPointer, len(name))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	importObj := wasmer.NewDefaultWasiImportObject()
	imports, err := importObj.Imports()
	imports, _ = imports.AppendFunction("GetMeasurementInt", GetMeasurementInt, C.GetMeasurementInt)
	imports, _ = imports.AppendFunction("GetMeasurementBytes", GetMeasurementBytes, C.GetMeasurementBytes)
	imports, _ = imports.AppendFunction("ForwardInterest", ForwardInterest, C.ForwardInterest)
	imports, _ = imports.AppendFunction("SetMeasurementInt", SetMeasurementInt, C.SetMeasurementInt)
	imports, _ = imports.AppendFunction("AddToMeasurementInt", AddToMeasurementInt, C.AddToMeasurementInt)
	imports, _ = imports.AppendFunction("SetMeasurementBytes", SetMeasurementBytes, C.SetMeasurementBytes)

	// Reads the WebAssembly module as bytes.
	bytes, _ := wasmer.ReadBytes("./csrc/sample.wasm")

	// Instantiates the WebAssembly module.
	instance, err = wasmer.NewInstanceWithImports(bytes, imports)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer instance.Close()

	measure = make(map[string]interface{})
	measure["counter"] = 0
	measure["hash"] = []byte{0, 0, 0, 0}

	callAfterReceiveInterest("hello")
	callAfterReceiveInterest("world")

	fmt.Printf("counter = %v\n", measure["counter"])
	fmt.Printf("hash = %x\n", measure["hash"])
}
