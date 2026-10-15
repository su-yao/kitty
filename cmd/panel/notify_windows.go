//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

const mbIconInfo = 0x00000040

// notify shows a native message box. The panel runs without a console
// (-H windowsgui), so this is the only way to surface fatal errors.
func notify(title, msg string) {
	t, _ := syscall.UTF16PtrFromString(title)
	m, _ := syscall.UTF16PtrFromString(msg)
	syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		0,
		uintptr(unsafe.Pointer(m)),
		uintptr(unsafe.Pointer(t)),
		mbIconInfo,
	)
}
