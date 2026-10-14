//go:build windows

package send

import (
	"encoding/binary"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                         = syscall.NewLazyDLL("user32.dll")
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procEnumWindows                = user32.NewProc("EnumWindows")
	procGetWindowTextW             = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW       = user32.NewProc("GetWindowTextLengthW")
	procIsWindowVisible            = user32.NewProc("IsWindowVisible")
	procSetForegroundWindow        = user32.NewProc("SetForegroundWindow")
	procShowWindow                 = user32.NewProc("ShowWindow")
	procGetForegroundWindow        = user32.NewProc("GetForegroundWindow")
	procAttachThreadInput          = user32.NewProc("AttachThreadInput")
	procGetWindowThreadProcessId   = user32.NewProc("GetWindowThreadProcessId")
	procAllowSetForegroundWindow   = user32.NewProc("AllowSetForegroundWindow")
	procBringWindowToTop           = user32.NewProc("BringWindowToTop")
	procSetFocus                   = user32.NewProc("SetFocus")
	procSetWindowPos               = user32.NewProc("SetWindowPos")
	procSendInput                  = user32.NewProc("SendInput")
	procGetWindow                  = user32.NewProc("GetWindow")
	procOpenClipboard              = user32.NewProc("OpenClipboard")
	procCloseClipboard             = user32.NewProc("CloseClipboard")
	procEmptyClipboard             = user32.NewProc("EmptyClipboard")
	procSetClipboardData           = user32.NewProc("SetClipboardData")
	procGetCurrentThreadId         = kernel32.NewProc("GetCurrentThreadId")
	procOpenProcess                = kernel32.NewProc("OpenProcess")
	procCloseHandle                = kernel32.NewProc("CloseHandle")
	procQueryFullProcessImageNameW = kernel32.NewProc("QueryFullProcessImageNameW")
	procGlobalAlloc                = kernel32.NewProc("GlobalAlloc")
	procGlobalLock                 = kernel32.NewProc("GlobalLock")
	procGlobalUnlock               = kernel32.NewProc("GlobalUnlock")
)

const (
	SW_RESTORE            = 9
	KEYEVENTF_KEYUP       = 0x0002
	KEYEVENTF_SCANCODE    = 0x0008
	INPUT_KEYBOARD        = 1
	CF_UNICODETEXT        = 13
	GMEM_MOVEABLE         = 0x0002
	PROCESS_QUERY_LIMITED = 0x1000
	GW_OWNER              = 4
	HWND_TOPMOST          = ^uintptr(0)
	HWND_NOTOPMOST        = ^uintptr(1)
	SWP_NOSIZE            = 0x0001
	SWP_NOMOVE            = 0x0002
	SWP_SHOWWINDOW        = 0x0040
	SC_GRAVE              = 0x29
	SC_ENTER              = 0x1C
	SC_LCTRL              = 0x1D
	SC_LSHIFT             = 0x2A
	SC_V                  = 0x2F
	SC_2                  = 0x03
	SC_3                  = 0x04
	SC_1                  = 0x02
	SC_F12                = 0x58
	SC_LALT               = 0x38
)

var usScan = map[rune]uint16{
	'1': 0x02, '2': 0x03, '3': 0x04, '4': 0x05, '5': 0x06,
	'6': 0x07, '7': 0x08, '8': 0x09, '9': 0x0A, '0': 0x0B,
	'q': 0x10, 'w': 0x11, 'e': 0x12, 'r': 0x13, 't': 0x14,
	'y': 0x15, 'u': 0x16, 'i': 0x17, 'o': 0x18, 'p': 0x19,
	'a': 0x1E, 's': 0x1F, 'd': 0x20, 'f': 0x21, 'g': 0x22,
	'h': 0x23, 'j': 0x24, 'k': 0x25, 'l': 0x26,
	'z': 0x2C, 'x': 0x2D, 'c': 0x2E, 'v': 0x2F, 'b': 0x30,
	'n': 0x31, 'm': 0x32,
	'-': 0x0C, '=': 0x0D, '[': 0x1A, ']': 0x1B,
	';': 0x27, '\'': 0x28, '`': 0x29, '\\': 0x2B,
	',': 0x33, '.': 0x34, '/': 0x35, ' ': 0x39,
}

type Options struct {
	ConsoleKey string
	DelayMS    int
	SkipOpen   bool
}

type kbdInput struct {
	Type  uint32
	Pad   uint32
	Vk    uint16
	Scan  uint16
	Flags uint32
	Time  uint32
	Extra uintptr
	Tail  [8]byte
}

func Send(command string, opt Options) error {
	command = strings.TrimSpace(command)
	if command == "" {
		return errors.New("空命令")
	}
	hwnd, _, err := findStellaris()
	if err != nil {
		return err
	}
	_ = focus(hwnd)
	d := time.Duration(opt.DelayMS) * time.Millisecond
	if d < 150*time.Millisecond {
		d = 200 * time.Millisecond
	}
	time.Sleep(d)
	if !opt.SkipOpen {
		openConsole(opt.ConsoleKey)
		time.Sleep(d)
	}
	if err := typeASCII(command); err != nil {
		_ = pasteClipboard(command)
	}
	time.Sleep(60 * time.Millisecond)
	tapScan(SC_ENTER)
	return nil
}

func FindStellaris() (string, error) {
	_, title, err := findStellaris()
	return title, err
}

func OpenConsole(key string) error {
	hwnd, title, err := findStellaris()
	if err != nil {
		return err
	}
	if err := focus(hwnd); err != nil {
		return fmt.Errorf("无法激活窗口 %s: %w", title, err)
	}
	time.Sleep(200 * time.Millisecond)
	openConsole(key)
	return nil
}

func findStellaris() (uintptr, string, error) {
	var found uintptr
	var title string
	cb := syscall.NewCallback(func(hwnd uintptr, _ uintptr) uintptr {
		vis, _, _ := procIsWindowVisible.Call(hwnd)
		if vis == 0 {
			return 1
		}
		owner, _, _ := procGetWindow.Call(hwnd, GW_OWNER)
		if owner != 0 {
			return 1
		}
		var pid uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
		if pid == 0 {
			return 1
		}
		if processName(pid) != "stellaris.exe" {
			return 1
		}
		t := windowTitle(hwnd)
		if isPanelTitle(t) {
			return 1
		}
		found = hwnd
		title = t
		if title == "" {
			title = "stellaris.exe"
		}
		return 0
	})
	procEnumWindows.Call(cb, 0)
	if found == 0 {
		return 0, "", errors.New("未找到 stellaris.exe。请启动游戏（非铁人），并以管理员运行本面板")
	}
	return found, title, nil
}

func isPanelTitle(t string) bool {
	low := strings.ToLower(t)
	return strings.Contains(t, "控制台面板") || strings.Contains(low, "console panel")
}

func processName(pid uint32) string {
	h, _, _ := procOpenProcess.Call(PROCESS_QUERY_LIMITED, 0, uintptr(pid))
	if h == 0 {
		return ""
	}
	defer procCloseHandle.Call(h)
	buf := make([]uint16, 32768)
	n := uint32(len(buf))
	r, _, _ := procQueryFullProcessImageNameW.Call(h, 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&n)))
	if r == 0 {
		return ""
	}
	return strings.ToLower(filepath.Base(syscall.UTF16ToString(buf)))
}

func windowTitle(hwnd uintptr) string {
	n, _, _ := procGetWindowTextLengthW.Call(hwnd)
	if n == 0 {
		return ""
	}
	buf := make([]uint16, n+1)
	procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), n+1)
	return syscall.UTF16ToString(buf)
}

func focus(hwnd uintptr) error {
	procAllowSetForegroundWindow.Call(^uintptr(0))
	procShowWindow.Call(hwnd, SW_RESTORE)
	fg, _, _ := procGetForegroundWindow.Call()
	var pid uint32
	fgTid, _, _ := procGetWindowThreadProcessId.Call(fg, uintptr(unsafe.Pointer(&pid)))
	curTid, _, _ := procGetCurrentThreadId.Call()
	attached := fgTid != 0 && curTid != 0 && fgTid != curTid
	if attached {
		procAttachThreadInput.Call(curTid, fgTid, 1)
		defer procAttachThreadInput.Call(curTid, fgTid, 0)
	}
	flags := uintptr(SWP_NOMOVE | SWP_NOSIZE | SWP_SHOWWINDOW)
	procSetWindowPos.Call(hwnd, HWND_TOPMOST, 0, 0, 0, 0, flags)
	procSetWindowPos.Call(hwnd, HWND_NOTOPMOST, 0, 0, 0, 0, flags)
	procBringWindowToTop.Call(hwnd)
	procSetForegroundWindow.Call(hwnd)
	procSetFocus.Call(hwnd)
	time.Sleep(80 * time.Millisecond)
	fg2, _, _ := procGetForegroundWindow.Call()
	if fg2 == hwnd {
		return nil
	}
	return errors.New("窗口未能置前，请先点一下群星画面再点按钮")
}

func openConsole(key string) {
	k := strings.ToLower(strings.TrimSpace(key))
	switch k {
	case "oem3", "~", "`", "tilde", "":
		tapScan(SC_GRAVE)
	case "shift+2":
		tapComboScan(SC_LSHIFT, SC_2)
	case "shift+3":
		tapComboScan(SC_LSHIFT, SC_3)
	case "alt+21", "alt+2+1":
		tapComboScan(SC_LALT, SC_2)
		time.Sleep(30 * time.Millisecond)
		tapScan(SC_1)
	case "f12":
		tapScan(SC_F12)
	default:
		tapScan(SC_GRAVE)
	}
}

func sendScan(scan uint16, up bool) {
	flags := uint32(KEYEVENTF_SCANCODE)
	if up {
		flags |= KEYEVENTF_KEYUP
	}
	in := kbdInput{Type: INPUT_KEYBOARD, Scan: scan, Flags: flags}
	procSendInput.Call(1, uintptr(unsafe.Pointer(&in)), unsafe.Sizeof(in))
}

func tapScan(scan uint16) {
	sendScan(scan, false)
	time.Sleep(18 * time.Millisecond)
	sendScan(scan, true)
	time.Sleep(18 * time.Millisecond)
}

func tapComboScan(mod, key uint16) {
	sendScan(mod, false)
	time.Sleep(12 * time.Millisecond)
	sendScan(key, false)
	time.Sleep(18 * time.Millisecond)
	sendScan(key, true)
	sendScan(mod, true)
	time.Sleep(18 * time.Millisecond)
}

func typeASCII(s string) error {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			sc, ok := usScan[r-'A'+'a']
			if !ok {
				return fmt.Errorf("无法输入 %q", r)
			}
			tapComboScan(SC_LSHIFT, sc)
			continue
		}
		if r == '_' {
			tapComboScan(SC_LSHIFT, 0x0C)
			continue
		}
		sc, ok := usScan[r]
		if !ok {
			return fmt.Errorf("无法输入 %q", r)
		}
		tapScan(sc)
	}
	return nil
}

func pasteClipboard(s string) error {
	if err := setClipboard(s); err != nil {
		return err
	}
	tapComboScan(SC_LCTRL, SC_V)
	time.Sleep(80 * time.Millisecond)
	return nil
}

func ptrFromUintptr(p uintptr) unsafe.Pointer {
	var ptr unsafe.Pointer
	*(*uintptr)(unsafe.Pointer(&ptr)) = p
	return ptr
}

func setClipboard(s string) error {
	u16, err := syscall.UTF16FromString(s)
	if err != nil {
		return err
	}
	size := len(u16) * 2
	h, _, err := procGlobalAlloc.Call(GMEM_MOVEABLE, uintptr(size))
	if h == 0 {
		return err
	}
	p, _, err := procGlobalLock.Call(h)
	if p == 0 {
		return err
	}
	dst := unsafe.Slice((*byte)(ptrFromUintptr(p)), size)
	for i, v := range u16 {
		binary.LittleEndian.PutUint16(dst[i*2:], v)
	}
	procGlobalUnlock.Call(h)
	ok, _, _ := procOpenClipboard.Call(0)
	if ok == 0 {
		return errors.New("OpenClipboard failed")
	}
	defer procCloseClipboard.Call()
	procEmptyClipboard.Call()
	r, _, err := procSetClipboardData.Call(CF_UNICODETEXT, h)
	if r == 0 {
		return err
	}
	return nil
}
