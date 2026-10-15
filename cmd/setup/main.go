package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	appName      = "Stellaris Console Panel"
	appNameZH    = "群星控制台面板"
	appPublisher = "Stellaris Console Panel"
	appVersion   = "1.1.3"
	exeName      = "StellarisConsolePanel.exe"
	uninstallID  = "StellarisConsolePanel"
)

//go:embed all:payload
var payloadFS embed.FS

type installOpts struct {
	Dest      string
	StartMenu bool
	Desktop   bool
	Launch    bool
	Silent    bool
	Uninstall bool
}

func main() {
	initConsole()
	opts := parseArgs()
	if opts.Uninstall {
		elevateIfNeeded()
		if err := runUninstall(opts); err != nil {
			fail(err)
		}
		return
	}
	elevateIfNeeded()
	if err := ensurePayload(); err != nil {
		fail(err)
	}
	if opts.Silent {
		if opts.Dest == "" {
			opts.Dest = defaultDir()
		}
		opts.StartMenu = true
		opts.Desktop = false
		opts.Launch = false
		if err := runInstall(opts); err != nil {
			fail(err)
		}
		return
	}
	if err := runWizard(opts); err != nil {
		fail(err)
	}
}

func parseArgs() installOpts {
	opts := installOpts{StartMenu: true, Desktop: true, Launch: true}
	for i := 1; i < len(os.Args); i++ {
		a := os.Args[i]
		switch strings.ToLower(a) {
		case "/s", "/silent", "-s", "--silent":
			opts.Silent = true
		case "/uninstall", "--uninstall", "/u":
			opts.Uninstall = true
		default:
			low := strings.ToLower(a)
			if strings.HasPrefix(low, "/d=") || strings.HasPrefix(low, "/dir=") {
				opts.Dest = strings.Trim(a[strings.Index(a, "=")+1:], `"`)
			}
		}
	}
	return opts
}

func ensurePayload() error {
	_, err := payloadFS.ReadFile("payload/" + exeName)
	if err != nil {
		return fmt.Errorf("安装包未包含 %s，请先运行 bash scripts/build.sh", exeName)
	}
	return nil
}

func runWizard(opts installOpts) error {
	in := bufio.NewReader(os.Stdin)
	fmt.Println()
	fmt.Println("================================================")
	fmt.Println("  " + appNameZH + "  安装向导")
	fmt.Println("  " + appName + "  v" + appVersion)
	fmt.Println("  适配群星 3.0.0 - 4.4.6+")
	fmt.Println("================================================")
	fmt.Println()
	fmt.Println("  本程序将安装到本地，并创建开始菜单与卸载项。")
	fmt.Println("  安装完成后请以管理员身份运行主程序。")
	fmt.Println()
	fmt.Print("按回车继续，或输入 Q 取消: ")
	if quit(readLine(in)) {
		fmt.Println("已取消。")
		return nil
	}

	opts.Dest = defaultDir()
	fmt.Println()
	fmt.Println("[1/3] 安装目录")
	fmt.Printf("  默认: %s\n", opts.Dest)
	fmt.Print("  新目录（回车使用默认）: ")
	if s := strings.TrimSpace(readLine(in)); s != "" {
		opts.Dest = s
	}

	fmt.Println()
	fmt.Println("[2/3] 附加选项")
	opts.StartMenu = askYes(in, "  创建开始菜单快捷方式", true)
	opts.Desktop = askYes(in, "  创建桌面快捷方式", true)
	opts.Launch = askYes(in, "  安装完成后启动程序", true)

	fmt.Println()
	fmt.Println("[3/3] 确认安装")
	fmt.Printf("  目录: %s\n", opts.Dest)
	fmt.Printf("  开始菜单: %v    桌面: %v    安装后启动: %v\n", opts.StartMenu, opts.Desktop, opts.Launch)
	fmt.Print("按回车开始安装，或输入 Q 取消: ")
	if quit(readLine(in)) {
		fmt.Println("已取消。")
		return nil
	}
	fmt.Println()
	if err := runInstall(opts); err != nil {
		return err
	}
	fmt.Println()
	fmt.Println("安装完成。")
	fmt.Println("  主程序: " + filepath.Join(opts.Dest, exeName))
	fmt.Println("  请先启动群星（非铁人），再运行本面板。")
	if !opts.Launch {
		fmt.Print("按回车退出: ")
		_, _ = in.ReadString('\n')
	}
	return nil
}

func runInstall(opts installOpts) error {
	if opts.Dest == "" {
		opts.Dest = defaultDir()
	}
	fmt.Println("正在复制文件...")
	n, err := extractPayload(opts.Dest)
	if err != nil {
		return err
	}
	fmt.Printf("  已复制 %d 个文件\n", n)

	uninst := filepath.Join(opts.Dest, "Uninstall.exe")
	self, err := os.Executable()
	if err != nil {
		return err
	}
	if err := copyFile(self, uninst); err != nil {
		return err
	}

	meta := map[string]string{
		"install_dir": opts.Dest,
		"version":     appVersion,
		"installed":   time.Now().Format(time.RFC3339),
	}
	b, _ := json.MarshalIndent(meta, "", "  ")
	_ = os.WriteFile(filepath.Join(opts.Dest, "install.json"), b, 0644)
	_ = os.WriteFile(filepath.Join(opts.Dest, "UNINSTALL.txt"), []byte(
		"卸载方式：\r\n"+
			"1. 运行本目录 Uninstall.exe\r\n"+
			"2. 或打开「设置 - 应用」，卸载「"+appNameZH+"」\r\n"+
			"3. 或运行: Uninstall.exe /S\r\n",
	), 0644)

	exe := filepath.Join(opts.Dest, exeName)
	if runtime.GOOS == "windows" {
		if opts.StartMenu {
			fmt.Println("正在创建开始菜单快捷方式...")
			if err := createShortcut(startMenuDir(), appNameZH+".lnk", exe, opts.Dest, appNameZH); err != nil {
				fmt.Println("  开始菜单快捷方式失败:", err)
			}
			_ = createShortcut(startMenuDir(), "卸载 "+appNameZH+".lnk", uninst, opts.Dest, "卸载 "+appNameZH)
		}
		if opts.Desktop {
			fmt.Println("正在创建桌面快捷方式...")
			if err := createShortcut(desktopDir(), appNameZH+".lnk", exe, opts.Dest, appNameZH); err != nil {
				fmt.Println("  桌面快捷方式失败:", err)
			}
		}
		fmt.Println("正在写入卸载信息...")
		if err := writeUninstallKey(opts.Dest, uninst); err != nil {
			fmt.Println("  卸载项写入失败:", err)
		}
		if opts.Launch {
			fmt.Println("正在启动程序...")
			_ = launchAdmin(exe)
		}
	}
	return nil
}

func runUninstall(opts installOpts) error {
	dest := opts.Dest
	if dest == "" {
		if d, err := readInstallDir(); err == nil {
			dest = d
		}
	}
	if dest == "" {
		dest = defaultDir()
	}
	fmt.Println("正在卸载 " + appNameZH + " ...")
	fmt.Println("  目录: " + dest)
	if !opts.Silent {
		fmt.Print("确认卸载？输入 Y 继续: ")
		in := bufio.NewReader(os.Stdin)
		if strings.ToLower(strings.TrimSpace(readLine(in))) != "y" {
			fmt.Println("已取消。")
			return nil
		}
	}
	terminateApp()
	time.Sleep(300 * time.Millisecond)
	removeShortcut(startMenuDir(), appNameZH+".lnk")
	removeShortcut(startMenuDir(), "卸载 "+appNameZH+".lnk")
	removeShortcut(desktopDir(), appNameZH+".lnk")
	_ = deleteUninstallKey()
	_ = os.Remove(filepath.Join(dest, "install.json"))
	_ = os.Remove(filepath.Join(dest, "UNINSTALL.txt"))
	_ = os.Remove(filepath.Join(dest, exeName))
	_ = os.Remove(filepath.Join(dest, "README.md"))
	_ = tryRemoveDir(dest)
	selfDelete(dest)
	fmt.Println("卸载完成。")
	return nil
}

func extractPayload(dest string) (int, error) {
	if err := os.MkdirAll(dest, 0755); err != nil {
		return 0, err
	}
	n := 0
	err := fs.WalkDir(payloadFS, "payload", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel("payload", path)
		if rel == "." {
			return nil
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		b, err := payloadFS.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(target, b, 0755); err != nil {
			return err
		}
		n++
		return nil
	})
	return n, err
}

func defaultDir() string {
	if runtime.GOOS == "windows" {
		pf := os.Getenv("ProgramFiles")
		if pf == "" {
			pf = `C:\Program Files`
		}
		return filepath.Join(pf, uninstallID)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, uninstallID)
}

func startMenuDir() string {
	prog := os.Getenv("ProgramData")
	if prog == "" {
		prog = `C:\ProgramData`
	}
	return filepath.Join(prog, "Microsoft", "Windows", "Start Menu", "Programs", appNameZH)
}

func desktopDir() string {
	pub := os.Getenv("PUBLIC")
	if pub == "" {
		pub = `C:\Users\Public`
	}
	return filepath.Join(pub, "Desktop")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func readLine(in *bufio.Reader) string {
	s, _ := in.ReadString('\n')
	return strings.TrimSpace(s)
}

func quit(s string) bool {
	s = strings.ToLower(s)
	return s == "q" || s == "quit" || s == "n" || s == "no"
}

func askYes(in *bufio.Reader, q string, def bool) bool {
	hint := "Y/n"
	if !def {
		hint = "y/N"
	}
	fmt.Printf("%s [%s]: ", q, hint)
	s := strings.ToLower(strings.TrimSpace(readLine(in)))
	if s == "" {
		return def
	}
	return s == "y" || s == "yes"
}

func tryRemoveDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return os.Remove(dir)
	}
	return nil
}

func fail(err error) {
	fmt.Println("安装失败:", err)
	if runtime.GOOS == "windows" {
		fmt.Print("按回车退出: ")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
	os.Exit(1)
}
