//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

func createShortcut(dir, name, target, workDir, desc string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	lnk := filepath.Join(dir, name)
	ps := fmt.Sprintf(
		"$s=(New-Object -ComObject WScript.Shell).CreateShortcut(%s); $s.TargetPath=%s; $s.WorkingDirectory=%s; $s.Description=%s; $s.WindowStyle=1; $s.Save()",
		psQuote(lnk), psQuote(target), psQuote(workDir), psQuote(desc),
	)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", ps)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := cmd.CombinedOutput(); err != nil {
		if e2 := writeURLShortcut(lnk, target); e2 == nil {
			return nil
		}
		return fmt.Errorf("快捷方式失败: %v %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func psQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func writeURLShortcut(path, target string) error {
	urlPath := path
	if strings.HasSuffix(strings.ToLower(urlPath), ".lnk") {
		urlPath = urlPath[:len(urlPath)-4] + ".url"
	}
	t := strings.ReplaceAll(target, `\`, `/`)
	body := "[InternetShortcut]\r\nURL=file:///" + t + "\r\n"
	return os.WriteFile(urlPath, []byte(body), 0644)
}

func removeShortcut(dir, name string) {
	_ = os.Remove(filepath.Join(dir, name))
	if strings.HasSuffix(strings.ToLower(name), ".lnk") {
		_ = os.Remove(filepath.Join(dir, name[:len(name)-4]+".url"))
	}
	entries, err := os.ReadDir(dir)
	if err == nil && len(entries) == 0 {
		_ = os.Remove(dir)
	}
}

func writeUninstallKey(dest, uninst string) error {
	k, err := createRegKey(`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\` + uninstallID)
	if err != nil {
		return err
	}
	defer k.close()
	size := estimateInstallSize(dest)
	vals := []error{
		k.setString("DisplayName", appNameZH+" ("+appName+")"),
		k.setString("DisplayVersion", appVersion),
		k.setString("Publisher", appPublisher),
		k.setString("InstallLocation", dest),
		k.setString("UninstallString", `"`+uninst+`" /uninstall`),
		k.setString("QuietUninstallString", `"`+uninst+`" /uninstall /S`),
		k.setString("DisplayIcon", filepath.Join(dest, exeName)),
		k.setDWORD("NoModify", 1),
		k.setDWORD("NoRepair", 1),
		k.setDWORD("EstimatedSize", uint32(size)),
		k.setString("InstallDate", time.Now().Format("20060102")),
	}
	for _, err := range vals {
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteUninstallKey() error {
	return deleteRegKey(`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\` + uninstallID)
}

func terminateApp() {
	taskkill := filepath.Join(os.Getenv("SystemRoot"), "System32", "taskkill.exe")
	if _, err := os.Stat(taskkill); err != nil {
		taskkill = "taskkill.exe"
	}
	cmd := exec.Command(taskkill, "/F", "/IM", exeName)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = cmd.Run()
}

func selfDelete(dest string) {
	exe := filepath.Join(dest, uninstExeName)
	if _, err := os.Stat(exe); err != nil {
		return
	}
	shell := os.Getenv("COMSPEC")
	if shell == "" {
		shell = "cmd.exe"
	}
	// Non-recursive rmdir only removes an empty directory. A recursive delete
	// would wipe unrelated files when the user installed into a shared folder
	// such as D:\Games. The delay lets this process exit before its image is
	// deleted.
	script := "ping -n 3 127.0.0.1 >nul & del /f /q \"" + exe + "\" & rmdir \"" + dest + "\" 2>nul"
	cmd := exec.Command(shell, "/c", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = cmd.Start()
}

func readInstallDir() (string, error) {
	k, err := openRegKey(`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\` + uninstallID)
	if err != nil {
		return "", err
	}
	defer k.close()
	return k.getString("InstallLocation")
}

func launchAdmin(exe string) error {
	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(exe)
	r, _, err := syscall.NewLazyDLL("shell32.dll").NewProc("ShellExecuteW").Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		0, 0, 1,
	)
	if r <= 32 {
		return err
	}
	return nil
}

func estimateInstallSize(dir string) int {
	var n int64
	_ = filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			n += info.Size()
		}
		return nil
	})
	return int(n / 1024)
}

type regKey struct{ h uintptr }

func createRegKey(path string) (*regKey, error) {
	p, _ := syscall.UTF16PtrFromString(path)
	var h uintptr
	var disp uint32
	r, _, err := syscall.NewLazyDLL("advapi32.dll").NewProc("RegCreateKeyExW").Call(
		0x80000002, uintptr(unsafe.Pointer(p)), 0, 0, 0, 0xF003F, 0,
		uintptr(unsafe.Pointer(&h)), uintptr(unsafe.Pointer(&disp)),
	)
	if r != 0 {
		return nil, err
	}
	return &regKey{h: h}, nil
}

func openRegKey(path string) (*regKey, error) {
	p, _ := syscall.UTF16PtrFromString(path)
	var h uintptr
	r, _, err := syscall.NewLazyDLL("advapi32.dll").NewProc("RegOpenKeyExW").Call(
		0x80000002, uintptr(unsafe.Pointer(p)), 0, 0x20019, uintptr(unsafe.Pointer(&h)),
	)
	if r != 0 {
		return nil, err
	}
	return &regKey{h: h}, nil
}

func deleteRegKey(path string) error {
	p, _ := syscall.UTF16PtrFromString(path)
	r, _, err := syscall.NewLazyDLL("advapi32.dll").NewProc("RegDeleteKeyW").Call(
		0x80000002, uintptr(unsafe.Pointer(p)),
	)
	if r != 0 {
		return err
	}
	return nil
}

func (k *regKey) close() {
	syscall.NewLazyDLL("advapi32.dll").NewProc("RegCloseKey").Call(k.h)
}

func (k *regKey) setString(name, val string) error {
	n, _ := syscall.UTF16PtrFromString(name)
	v, _ := syscall.UTF16FromString(val)
	r, _, err := syscall.NewLazyDLL("advapi32.dll").NewProc("RegSetValueExW").Call(
		k.h, uintptr(unsafe.Pointer(n)), 0, 1, uintptr(unsafe.Pointer(&v[0])), uintptr(len(v)*2),
	)
	if r != 0 {
		return err
	}
	return nil
}

func (k *regKey) setDWORD(name string, val uint32) error {
	n, _ := syscall.UTF16PtrFromString(name)
	r, _, err := syscall.NewLazyDLL("advapi32.dll").NewProc("RegSetValueExW").Call(
		k.h, uintptr(unsafe.Pointer(n)), 0, 4, uintptr(unsafe.Pointer(&val)), 4,
	)
	if r != 0 {
		return err
	}
	return nil
}

func (k *regKey) getString(name string) (string, error) {
	n, _ := syscall.UTF16PtrFromString(name)
	var typ, size uint32
	proc := syscall.NewLazyDLL("advapi32.dll").NewProc("RegQueryValueExW")
	r, _, err := proc.Call(k.h, uintptr(unsafe.Pointer(n)), 0, uintptr(unsafe.Pointer(&typ)), 0, uintptr(unsafe.Pointer(&size)))
	if r != 0 {
		return "", err
	}
	buf := make([]uint16, size/2+1)
	r, _, err = proc.Call(k.h, uintptr(unsafe.Pointer(n)), 0, uintptr(unsafe.Pointer(&typ)), uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if r != 0 {
		return "", err
	}
	return syscall.UTF16ToString(buf), nil
}
