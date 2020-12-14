(module
  ;; import functions
  (import "env" "AddToMeasurementInt" (func $AddToMeasurementInt (param (;name;) i32 (;value;) i32)))
  (import "env" "GetMeasurementBytes" (func $GetMeasurementBytes (param (;name;) i32 (;buf;) i32 (;buflen;) i32) (result i32)))
  (import "env" "SetMeasurementBytes" (func $SetMeasurementBytes (param (;name;) i32 (;value;) i32 (;size;) i32)))
  (import "env" "ForwardInterest" (func $ForwardInterest (param (;name;) i32)))
  ;; declear memory (1 page) and export
  (memory (export "memory") 1)
  ;; global variables
  (global $memPtr (mut i32) (i32.const 0)) ;; Cyclic memory ptr
  (global $memLimit i32 (i32.const 256))

  ;; allocate memory cyclicly, without freeing any block
  (func $Allocate (param $size i32) (result i32)
    (local $ret i32)
    (local.set $ret (global.get $memPtr))
    (global.set $memPtr
      (i32.add
        (global.get $memPtr)
        (local.get $size)))
    (if (i32.ge_u (global.get $memPtr) (global.get $memLimit))
      (global.set $memPtr (i32.const 0)))
    (local.get $ret)
  )

  (func $Deallocate (param $ptr i32)
    nop)

  ;; after receive interest
  (func (export "AfterReceiveInterest") (param $name i32) (param $nameLen i32)
    (local $buf i32)
    (local $i i32)
    (call $AddToMeasurementInt (i32.const 261) (i32.const 1))
    (local.set $buf (call $Allocate (i32.const 4)))
    (drop (call $GetMeasurementBytes (i32.const 256) (local.get $buf) (i32.const 4)))
    (local.set $i (i32.const 0))
    (block
      (loop
        (i32.store8 (i32.add (local.get $buf) (local.get $i))
          (i32.xor
            (i32.load8_u (i32.add (local.get $buf) (local.get $i)))
            (i32.load8_u (i32.add (local.get $name) (local.get $i)))))
        (local.set $i (i32.add (local.get $i) (i32.const 1)))
        (br_if 1 (i32.ge_u (local.get $i) (i32.const 4)))
        (br 0)))
    (call $SetMeasurementBytes (i32.const 256) (local.get $buf) (i32.const 4))
    (call $ForwardInterest (i32.load8_u (local.get $name)))
  )

  (export "Allocate" (func $Allocate))
  (export "Deallocate" (func $Deallocate))
  (data (i32.const 256) "hash\00counter\00")
)
