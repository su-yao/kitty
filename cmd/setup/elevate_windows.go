//go:build windows

package main

import (
	"fmt"
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
	cwd, _ := os.Getwd()
	dir, _ := syscall.UTF16PtrFromString(cwd)
	var args *uint16
	if len(os.Args) > 1 {
		a, _ := syscall.UTF16PtrFromString(quoteArgs(os.Args[1:]))
		args = a
	}
	r, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("ShellExecuteW").Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		uintptr(unsafe.Pointer(args)),
		uintptr(unsafe.Pointer(dir)),
		1,
	)
	if r > 32 {
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, "未能获取管理员权限。安装到 Program Files 可能失败。")
}

func alreadyAdmin() bool {
	r, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("IsUserAnAdmin").Call()
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
