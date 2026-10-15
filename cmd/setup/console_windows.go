//go:build windows

package main

import "syscall"

// initConsole switches the console code page to UTF-8 so the Chinese wizard
// text renders correctly on Windows 10 instead of showing mojibake.
func initConsole() {
	k := syscall.NewLazyDLL("kernel32.dll")
	k.NewProc("SetConsoleOutputCP").Call(65001)
	k.NewProc("SetConsoleCP").Call(65001)
}
