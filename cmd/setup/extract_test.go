package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractPayload(t *testing.T) {
	dir := t.TempDir()
	n, err := extractPayload(dir)
	if err != nil {
		t.Fatal(err)
	}
	if n < 1 {
		t.Fatalf("expected files, got %d", n)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 1 {
		t.Fatalf("expected extracted content in %s", dir)
	}
}

func TestParseArgsSilent(t *testing.T) {
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"setup", "/S", `/DIR=D:\Apps\Panel`}
	opts := parseArgs()
	if !opts.Silent {
		t.Fatal("silent")
	}
	if opts.Dest != `D:\Apps\Panel` {
		t.Fatalf("dir=%q", opts.Dest)
	}
}

func TestParseArgsQuotedDir(t *testing.T) {
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"setup", `/DIR="C:\Program Files\Panel"`}
	opts := parseArgs()
	if opts.Dest != `C:\Program Files\Panel` {
		t.Fatalf("dir=%q", opts.Dest)
	}
}

func TestSanitizeFileName(t *testing.T) {
	ok := map[string]string{
		"StellarisConsolePanel.exe": "StellarisConsolePanel.exe",
		"README.md":                 "README.md",
		"Uninstall.txt":             "Uninstall.txt",
		"a b.txt":                   "ab.txt",
	}
	for in, want := range ok {
		if got := sanitizeFileName(in); got != want {
			t.Fatalf("%q -> %q, want %q", in, got, want)
		}
	}
	bad := []string{"CON", "nul.txt", "..", ".", "", "???", "卸载.lnk"}
	for _, in := range bad {
		if got := sanitizeFileName(in); got != "" {
			t.Fatalf("want drop %q, got %q", in, got)
		}
	}
}

func TestSafeRelPath(t *testing.T) {
	got, ok := safeRelPath(`payload\..\StellarisConsolePanel.exe`)
	if ok {
		t.Fatalf("parent path should be rejected: %q", got)
	}
	got, ok = safeRelPath("StellarisConsolePanel.exe")
	if !ok || got != "StellarisConsolePanel.exe" {
		t.Fatalf("exe path: %q %v", got, ok)
	}
}

func TestWriteFileRetry(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.bin")
	if err := writeFileRetry(p, []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "ok" {
		t.Fatalf("got %q", b)
	}
	if err := writeFileRetry(p, []byte("two"), 0644); err != nil {
		t.Fatal(err)
	}
	b, err = os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "two" {
		t.Fatalf("overwrite got %q", b)
	}
}

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.bin")
	dst := filepath.Join(dir, "dst.bin")
	if err := os.WriteFile(src, []byte("payload"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := copyFile(src, dst); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "payload" {
		t.Fatalf("got %q", b)
	}
}

func TestVersionPart(t *testing.T) {
	if versionPart("1.4.6", 0) != 1 || versionPart("1.4.6", 1) != 4 || versionPart("1.4.6", 2) != 6 {
		t.Fatalf("versionPart 1.4.6 -> %d.%d.%d", versionPart("1.4.6", 0), versionPart("1.4.6", 1), versionPart("1.4.6", 2))
	}
	if versionPart("1.4.6", 3) != 0 || versionPart("", 0) != 0 {
		t.Fatal("versionPart bounds")
	}
}

func TestParseArgsUninstall(t *testing.T) {
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"setup", "/uninstall", "/S"}
	opts := parseArgs()
	if !opts.Uninstall || !opts.Silent {
		t.Fatalf("%+v", opts)
	}
}
