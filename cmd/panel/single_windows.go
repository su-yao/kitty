//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

const errorAlreadyExists = 183

// singleInstance keeps a second copy of the panel from starting a second
// server and injecting keystrokes at the same time.
func singleInstance() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	name, _ := syscall.UTF16PtrFromString(`Local\StellarisConsolePanel`)
	h, _, err := kernel32.NewProc("CreateMutexW").Call(0, 0, uintptr(unsafe.Pointer(name)))
	if errno, ok := err.(syscall.Errno); ok && errno == errorAlreadyExists {
		if h != 0 {
			kernel32.NewProc("CloseHandle").Call(h)
		}
		return false
	}
	return true
}
