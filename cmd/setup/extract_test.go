package main

import (
	"os"
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

func TestParseArgsUninstall(t *testing.T) {
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"setup", "/uninstall", "/S"}
	opts := parseArgs()
	if !opts.Uninstall || !opts.Silent {
		t.Fatalf("%+v", opts)
	}
}
