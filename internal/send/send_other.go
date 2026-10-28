//go:build !windows

package send

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Options struct {
	ConsoleKey string
	DelayMS    int
	SkipOpen   bool
}

var (
	mu      sync.Mutex
	history []string
)

func SendChain(commands []string, opt Options) error {
	n := 0
	var last error
	for _, c := range commands {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		last = Send(c, opt)
		n++
	}
	if n == 0 {
		return errors.New("空命令")
	}
	return last
}

func Send(command string, opt Options) error {
	if command == "" {
		return errors.New("空命令")
	}
	mu.Lock()
	history = append(history, fmt.Sprintf("%s\t%s", time.Now().Format(time.RFC3339), command))
	mu.Unlock()
	dir, _ := os.UserConfigDir()
	if dir == "" {
		dir = os.TempDir()
	}
	p := filepath.Join(dir, "stellaris-console-panel", "sent.log")
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	f, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("预览模式已记录命令，但写日志失败: %w", err)
	}
	defer f.Close()
	_, _ = fmt.Fprintf(f, "%s\tkey=%s\t%s\n", time.Now().Format(time.RFC3339), opt.ConsoleKey, command)
	return errors.New("当前不是 Windows：命令已写入日志，未发送到群星。请在 Win10 上运行安装包")
}

func FindStellaris() (string, error) {
	return "", errors.New("仅 Windows 可检测群星窗口")
}

func OpenConsole(string) error {
	return errors.New("仅 Windows 可打开群星控制台")
}
