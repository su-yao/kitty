//go:build windows

package main

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func elevateIfNeeded() {
	if alreadyAdmin() {
		return
	}
	exe, err := os.Executable()
	if err != nil {
		return
	}
	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(exe)
	dir, _ := syscall.UTF16PtrFromString("")
	var args *uint16
	if len(os.Args) > 1 {
		a, _ := syscall.UTF16PtrFromString(quoteArgs(os.Args[1:]))
		args = a
	}
	shell32 := syscall.NewLazyDLL("shell32.dll")
	proc := shell32.NewProc("ShellExecuteW")
	r, _, _ := proc.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(file)), uintptr(unsafe.Pointer(args)), uintptr(unsafe.Pointer(dir)), 1)
	if r > 32 {
		os.Exit(0)
	}
}

func alreadyAdmin() bool {
	mod := syscall.NewLazyDLL("shell32.dll")
	proc := mod.NewProc("IsUserAnAdmin")
	r, _, _ := proc.Call()
	return r != 0
}

func quoteArgs(args []string) string {
	parts := make([]string, len(args))
	for i, a := range args {
		if strings.ContainsAny(a, " \t") {
			parts[i] = `"` + strings.ReplaceAll(a, `"`, `\"`) + `"`
			continue
		}
		parts[i] = a
	}
	return strings.Join(parts, " ")
}
