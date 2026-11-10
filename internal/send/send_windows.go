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
	procGlobalFree                 = kernel32.NewProc("GlobalFree")
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
	SC_F12                = 0x58
	SC_LALT               = 0x38
	SC_C                  = 0x2E
)

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

func SendChain(commands []string, opt Options) error {
	first := true
	n := 0
	for _, c := range commands {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		o := opt
		if !first {
			o.SkipOpen = true
		}
		if err := Send(c, o); err != nil {
			return err
		}
		first = false
		n++
		time.Sleep(120 * time.Millisecond)
	}
	if n == 0 {
		return errors.New("空命令")
	}
	return nil
}

func Send(command string, opt Options) error {
	command = strings.TrimSpace(command)
	if command == "" {
		return errors.New("空命令")
	}
	hwnd, title, err := findStellaris()
	if err != nil {
		return err
	}
	if err := focus(hwnd); err != nil {
		return fmt.Errorf("无法激活窗口 %s：%w", title, err)
	}
	d := time.Duration(opt.DelayMS) * time.Millisecond
	if d < 150*time.Millisecond {
		d = 200 * time.Millisecond
	}
	time.Sleep(d)
	if !opt.SkipOpen {
		if err := openConsole(opt.ConsoleKey); err != nil {
			return err
		}
		time.Sleep(d)
	}
	// Decide the transport before typing anything. Switching to the clipboard
	// after a partial type would leave the already typed prefix in the console
	// and append the whole command again.
	if canTypeASCII(command) {
		if err := typeASCII(command); err != nil {
			return err
		}
	} else {
		if err := pasteClipboard(command); err != nil {
			return fmt.Errorf("无法输入该命令，且写入剪贴板失败: %w", err)
		}
	}
	time.Sleep(60 * time.Millisecond)
	return tapScan(SC_ENTER)
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
	return openConsole(key)
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

func openConsole(key string) error {
	k := strings.ToLower(strings.TrimSpace(key))
	switch k {
	case "oem3", "~", "`", "tilde", "^", "caret", "°", "degree", "²", "squared", "":
		return tapScan(SC_GRAVE)
	case "shift+2":
		return tapComboScan(SC_LSHIFT, SC_2)
	case "shift+3", "§", "section":
		return tapComboScan(SC_LSHIFT, SC_3)
	case "alt+21", "alt+2+1":
		return tapAltNumpad(VK_NUMPAD2, VK_NUMPAD1)
	case "shift+alt+c":
		return tapComboScanMods(SC_LSHIFT, SC_LALT, SC_C)
	case "f12":
		return tapScan(SC_F12)
	default:
		return tapScan(SC_GRAVE)
	}
}

func sendInput(in kbdInput) error {
	n, _, err := procSendInput.Call(1, uintptr(unsafe.Pointer(&in)), unsafe.Sizeof(in))
	if n != 1 {
		if err == nil || err == syscall.Errno(0) {
			return errors.New("SendInput 失败")
		}
		return err
	}
	return nil
}

func sendScan(scan uint16, up bool) error {
	flags := uint32(KEYEVENTF_SCANCODE)
	if up {
		flags |= KEYEVENTF_KEYUP
	}
	return sendInput(kbdInput{Type: INPUT_KEYBOARD, Scan: scan, Flags: flags})
}

func sendVK(vk uint16, up bool) error {
	flags := uint32(0)
	if up {
		flags = KEYEVENTF_KEYUP
	}
	return sendInput(kbdInput{Type: INPUT_KEYBOARD, Vk: vk, Flags: flags})
}

func tapScan(scan uint16) error {
	if err := sendScan(scan, false); err != nil {
		return err
	}
	time.Sleep(18 * time.Millisecond)
	if err := sendScan(scan, true); err != nil {
		return err
	}
	time.Sleep(18 * time.Millisecond)
	return nil
}

func tapAltNumpad(keys ...uint16) error {
	if err := sendVK(VK_MENU, false); err != nil {
		return err
	}
	time.Sleep(12 * time.Millisecond)
	for _, vk := range keys {
		if err := sendVK(vk, false); err != nil {
			_ = sendVK(VK_MENU, true)
			return err
		}
		time.Sleep(18 * time.Millisecond)
		if err := sendVK(vk, true); err != nil {
			_ = sendVK(VK_MENU, true)
			return err
		}
		time.Sleep(12 * time.Millisecond)
	}
	if err := sendVK(VK_MENU, true); err != nil {
		return err
	}
	time.Sleep(18 * time.Millisecond)
	return nil
}

func tapComboScanMods(mod1, mod2, key uint16) error {
	if err := sendScan(mod1, false); err != nil {
		return err
	}
	time.Sleep(12 * time.Millisecond)
	if err := sendScan(mod2, false); err != nil {
		return err
	}
	time.Sleep(12 * time.Millisecond)
	if err := sendScan(key, false); err != nil {
		return err
	}
	time.Sleep(18 * time.Millisecond)
	if err := sendScan(key, true); err != nil {
		return err
	}
	if err := sendScan(mod2, true); err != nil {
		return err
	}
	if err := sendScan(mod1, true); err != nil {
		return err
	}
	time.Sleep(18 * time.Millisecond)
	return nil
}

func tapComboScan(mod, key uint16) error {
	if err := sendScan(mod, false); err != nil {
		return err
	}
	time.Sleep(12 * time.Millisecond)
	if err := sendScan(key, false); err != nil {
		return err
	}
	time.Sleep(18 * time.Millisecond)
	if err := sendScan(key, true); err != nil {
		return err
	}
	if err := sendScan(mod, true); err != nil {
		return err
	}
	time.Sleep(18 * time.Millisecond)
	return nil
}

func typeASCII(s string) error {
	for _, r := range s {
		if sc, ok := usScanShift[r]; ok {
			if err := tapComboScan(SC_LSHIFT, sc); err != nil {
				return err
			}
			continue
		}
		if r >= 'A' && r <= 'Z' {
			sc, ok := usScan[r-'A'+'a']
			if !ok {
				return fmt.Errorf("无法输入 %q", r)
			}
			if err := tapComboScan(SC_LSHIFT, sc); err != nil {
				return err
			}
			continue
		}
		sc, ok := usScan[r]
		if !ok {
			return fmt.Errorf("无法输入 %q", r)
		}
		if err := tapScan(sc); err != nil {
			return err
		}
	}
	return nil
}

func pasteClipboard(s string) error {
	if err := setClipboard(s); err != nil {
		return err
	}
	if err := tapComboScan(SC_LCTRL, SC_V); err != nil {
		return err
	}
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
	// The clipboard owns the handle only after a successful SetClipboardData.
	owned := false
	defer func() {
		if !owned {
			procGlobalFree.Call(h)
		}
	}()
	p, _, err := procGlobalLock.Call(h)
	if p == 0 {
		return err
	}
	dst := unsafe.Slice((*byte)(ptrFromUintptr(p)), size)
	for i, v := range u16 {
		binary.LittleEndian.PutUint16(dst[i*2:], v)
	}
	procGlobalUnlock.Call(h)
	opened := false
	var openErr error
	for i := 0; i < 8; i++ {
		ok, _, err := procOpenClipboard.Call(0)
		if ok != 0 {
			opened = true
			break
		}
		openErr = err
		time.Sleep(25 * time.Millisecond)
	}
	if !opened {
		if openErr == nil || openErr == syscall.Errno(0) {
			return errors.New("OpenClipboard failed")
		}
		return openErr
	}
	defer procCloseClipboard.Call()
	procEmptyClipboard.Call()
	r, _, err := procSetClipboardData.Call(CF_UNICODETEXT, h)
	if r == 0 {
		return err
	}
	owned = true
	return nil
}
