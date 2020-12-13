// If any standard library is used, one must use wasic++ to compile it.
#include <cstdlib>

// Expose malloc and free for the host to allocate memory
// Since WASM uses linear memory, everything must be copied between Host memory and WASM.
void* Allocate(int size) asm("Allocate");
void Deallocate(void* ptr) asm("Deallocate");

// The following are just proposals to illustrate how to interact with the host.
// We need to discuss what interface we want to use.
extern int GetMeasurementInt(const char* name) asm("GetMeasurementInt");
// We can define C++ structures, but only their pointers can be used in export/import functions.
// We need to copy data in and out. This interface is just for example. I think we can do better.
extern int GetMeasurementBytes(const char* name, void* buf, int buflen) asm("GetMeasurementBytes");
extern void SetMeasurementInt(const char* name, int value) asm("SetMeasurementInt");
extern void AddToMeasurementInt(const char* name, int value) asm("AddToMeasurementInt");
extern void SetMeasurementBytes(const char* name, const void* value, int size) asm("SetMeasurementBytes");
extern void ForwardInterest(int faceId) asm("ForwardInterest");

void AfterReceiveInterest(unsigned char* name, int nameLen) asm("AfterReceiveInterest");

void* Allocate(int size) {
  return malloc(size);
}

void Deallocate(void* ptr) {
  free(ptr);
}

void AfterReceiveInterest(unsigned char* name, int nameLen) {
  unsigned char buf[4];
  int i;

  AddToMeasurementInt("counter", 1);
  GetMeasurementBytes("hash", buf, 4);
  for(i = 0; i < 4; i ++){
    buf[i] ^= name[i];
  }
  SetMeasurementBytes("hash", buf, 4);
  ForwardInterest(int(name[0]));
}

int main() {
  return 0;
}
