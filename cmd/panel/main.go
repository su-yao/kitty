package main

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"stellaris-console-panel/internal/catalog"
	"stellaris-console-panel/internal/scan"
	"stellaris-console-panel/internal/send"
)

//go:embed all:web
var webFS embed.FS

// appVersion is injected at build time via -ldflags "-X main.appVersion=...".
var appVersion = "dev"

type Settings struct {
	GameRoot   string `json:"game_root"`
	Workshop   string `json:"workshop"`
	ExtraMods  string `json:"extra_mods"`
	ConsoleKey string `json:"console_key"`
	Version    string `json:"version"`
}

type App struct {
	store *catalog.Store
	mu    sync.Mutex
	set   Settings
	// sendMu serializes keystroke injection. Two concurrent requests would
	// interleave their keystrokes and corrupt both commands in the console.
	sendMu sync.Mutex
}

func main() {
	elevateIfNeeded()
	if !singleInstance() {
		notify("群星控制台面板", "面板已经在运行，请使用已打开的窗口。")
		return
	}
	app := &App{store: catalog.New(), set: normalizeSettings(loadSettings())}
	mux := http.NewServeMux()
	static, err := fs.Sub(webFS, "web")
	if err != nil {
		die(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		b, _ := webFS.ReadFile("web/index.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
	})
	mux.HandleFunc("/api/meta", app.meta)
	mux.HandleFunc("/api/commands", app.commands)
	mux.HandleFunc("/api/send", app.sendCmd)
	mux.HandleFunc("/api/open_console", app.openConsole)
	mux.HandleFunc("/api/scan", app.scanNow)
	mux.HandleFunc("/api/detect", app.detect)
	mux.HandleFunc("/api/settings", app.settings)
	addr := "127.0.0.1:0"
	if p := os.Getenv("PORT"); p != "" {
		addr = "127.0.0.1:" + p
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		die(err)
	}
	url := "http://" + ln.Addr().String() + "/"
	go func() {
		time.Sleep(300 * time.Millisecond)
		openUI(url)
	}()
	log.Println("Stellaris Console Panel", url)
	if err := http.Serve(ln, secureHeaders(localGuard(mux))); err != nil {
		die(err)
	}
}

func die(err error) {
	notify("群星控制台面板", err.Error())
	log.Fatal(err)
}

// previewHostSuffix is the trusted reverse-proxy domain used by the online
// preview. Requests arriving through the proxy carry this domain in Host.
const previewHostSuffix = ".monkeycode-ai.online"

// localGuard rejects requests whose Host or Origin is neither loopback nor the
// trusted preview domain. This blocks DNS-rebinding and cross-site requests
// from reaching the keystroke injection API while still allowing the preview
// proxy to forward requests.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self'; img-src 'self' data:; frame-ancestors 'none'")
		next.ServeHTTP(w, r)
	})
}

func localGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !allowedHost(r.Host) {
			httpError(w, http.StatusForbidden, "invalid host")
			return
		}
		if origin := r.Header.Get("Origin"); origin != "" {
			u, err := url.Parse(origin)
			if err != nil || !allowedHost(u.Host) {
				httpError(w, http.StatusForbidden, "invalid origin")
				return
			}
		}
		// JSON-body endpoints accept JSON only. A cross-site HTML form cannot
		// produce that content type, which closes the remaining CSRF gap for
		// requests that arrive without an Origin header.
		if r.Method == http.MethodPost && needsJSONBody(r.URL.Path) &&
			!strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			httpError(w, http.StatusUnsupportedMediaType, "content type must be application/json")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func needsJSONBody(path string) bool {
	switch path {
	case "/api/send", "/api/open_console", "/api/settings", "/api/scan":
		return true
	}
	return false
}

func allowedHost(hostport string) bool {
	if hostport == "" {
		return false
	}
	host := hostport
	if h, _, err := net.SplitHostPort(hostport); err == nil {
		host = h
	}
	host = strings.ToLower(strings.Trim(host, "[]"))
	if host == "127.0.0.1" || host == "localhost" || host == "::1" {
		return true
	}
	return strings.HasSuffix(host, previewHostSuffix)
}

func (a *App) snapSettings() Settings {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.set
}

func (a *App) meta(w http.ResponseWriter, _ *http.Request) {
	s := a.snapSettings()
	writeJSON(w, map[string]any{
		"version":         appVersion,
		"categories":      catalog.Categories,
		"versions":        catalog.Versions,
		"default_version": or(s.Version, "4.4.6"),
		"settings":        s,
		"scanned":         a.store.ScannedCount(),
	})
}

func (a *App) commands(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	size, _ := strconv.Atoi(q.Get("page_size"))
	p := a.store.List(q.Get("category"), q.Get("q"), q.Get("version"), page, size)
	writeJSON(w, map[string]any{
		"items":       p.Items,
		"total":       p.Total,
		"page":        p.Page,
		"page_size":   p.PageSize,
		"total_pages": p.TotalPages,
		"category":    p.Category,
		"scanned":     a.store.ScannedCount(),
	})
}

func (a *App) sendCmd(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Command    string `json:"command"`
		ConsoleKey string `json:"console_key"`
		SkipOpen   bool   `json:"skip_open"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpError(w, 400, err.Error())
		return
	}
	if err := validateCommand(body.Command); err != nil {
		httpError(w, 400, err.Error())
		return
	}
	key := body.ConsoleKey
	if key == "" {
		key = a.snapSettings().ConsoleKey
	}
	if !validConsoleKey(key) {
		key = "oem3"
	}
	a.sendMu.Lock()
	err := send.Send(body.Command, send.Options{ConsoleKey: key, DelayMS: 180, SkipOpen: body.SkipOpen})
	a.sendMu.Unlock()
	if err != nil {
		httpError(w, 500, err.Error())
		return
	}
	writeJSON(w, map[string]any{"ok": true, "message": "已发送: " + body.Command})
}

func validateCommand(cmd string) error {
	c := strings.TrimSpace(cmd)
	if c == "" {
		return errors.New("空命令")
	}
	if len(c) > 512 {
		return errors.New("命令过长")
	}
	if strings.ContainsAny(c, "\r\n") {
		return errors.New("命令不能包含换行")
	}
	return nil
}

func (a *App) openConsole(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ConsoleKey string `json:"console_key"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&body)
	}
	key := body.ConsoleKey
	if key == "" {
		key = a.snapSettings().ConsoleKey
	}
	if !validConsoleKey(key) {
		key = "oem3"
	}
	a.sendMu.Lock()
	err := send.OpenConsole(key)
	a.sendMu.Unlock()
	if err != nil {
		httpError(w, 500, err.Error())
		return
	}
	writeJSON(w, map[string]any{"ok": true, "message": "已尝试打开群星控制台"})
}

func (a *App) scanNow(w http.ResponseWriter, _ *http.Request) {
	s := a.snapSettings()
	if s.GameRoot == "" && s.Workshop == "" && s.ExtraMods == "" {
		httpError(w, 400, "请至少填写游戏根目录或 MOD 目录")
		return
	}
	ver := scan.GameVersion(s.GameRoot)
	if ver == "" {
		ver = s.Version
	}
	items, err := scan.Scan(s.GameRoot, s.Workshop, s.ExtraMods, ver)
	if err != nil {
		httpError(w, 500, err.Error())
		return
	}
	a.store.ReplaceScanned(items)
	writeJSON(w, map[string]any{
		"ok":      true,
		"count":   len(items),
		"version": ver,
		"message": "扫描完成 " + strconv.Itoa(len(items)) + " 条，检测到版本 " + or(ver, "未知"),
	})
}

func (a *App) detect(w http.ResponseWriter, _ *http.Request) {
	title, err := send.FindStellaris()
	if err != nil {
		httpError(w, 404, err.Error())
		return
	}
	writeJSON(w, map[string]any{"ok": true, "message": "找到窗口: " + title})
}

func (a *App) settings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		writeJSON(w, a.snapSettings())
		return
	}
	// Decode on top of the current values so omitted fields (for example
	// extra_mods from an older client) are preserved instead of being wiped.
	s := a.snapSettings()
	r.Body = http.MaxBytesReader(w, r.Body, 32<<10)
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		httpError(w, 400, err.Error())
		return
	}
	s = normalizeSettings(s)
	if err := saveSettings(s); err != nil {
		httpError(w, 500, "保存设置失败: "+err.Error())
		return
	}
	a.mu.Lock()
	a.set = s
	a.mu.Unlock()
	writeJSON(w, map[string]any{"ok": true})
}

func normalizeSettings(s Settings) Settings {
	trim := func(v string) string {
		return strings.Trim(strings.TrimSpace(v), `"`)
	}
	s.GameRoot = trim(s.GameRoot)
	s.Workshop = trim(s.Workshop)
	s.ExtraMods = trim(s.ExtraMods)
	if !validVersion(s.Version) {
		s.Version = "4.4.6"
	}
	if !validConsoleKey(s.ConsoleKey) {
		s.ConsoleKey = "oem3"
	}
	return s
}

func validVersion(v string) bool {
	for _, item := range catalog.Versions {
		if item == v {
			return true
		}
	}
	return false
}

func validConsoleKey(k string) bool {
	switch strings.ToLower(strings.TrimSpace(k)) {
	case "oem3", "~", "`", "tilde", "shift+2", "shift+3", "alt+21", "alt+2+1", "f12":
		return true
	}
	return false
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func configPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	p := filepath.Join(dir, "StellarisConsolePanel")
	_ = os.MkdirAll(p, 0755)
	return filepath.Join(p, "settings.json")
}

func loadSettings() Settings {
	s := Settings{Version: "4.4.6", ConsoleKey: "oem3"}
	b, err := os.ReadFile(configPath())
	if err != nil {
		return s
	}
	_ = json.Unmarshal(b, &s)
	if s.Version == "" {
		s.Version = "4.4.6"
	}
	if s.ConsoleKey == "" {
		s.ConsoleKey = "oem3"
	}
	return s
}

func saveSettings(s Settings) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	path := configPath()
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func or(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

func openUI(url string) {
	if runtime.GOOS == "windows" {
		msedge := filepath.Join(os.Getenv("ProgramFiles(x86)"), "Microsoft", "Edge", "Application", "msedge.exe")
		if _, err := os.Stat(msedge); err == nil {
			_ = exec.Command(msedge, "--app="+url, "--start-maximized", "--window-size=1600,900").Start()
			return
		}
		chrome := filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe")
		if _, err := os.Stat(chrome); err == nil {
			_ = exec.Command(chrome, "--app="+url, "--start-maximized", "--window-size=1600,900").Start()
			return
		}
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		return
	}
	if _, err := exec.LookPath("xdg-open"); err == nil {
		_ = exec.Command("xdg-open", url).Start()
		return
	}
	log.Println("打开浏览器访问", url)
	notify("群星控制台面板", "未找到浏览器，请手动打开：\n"+url)
}
