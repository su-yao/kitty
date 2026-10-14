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

func TestParseArgsUninstall(t *testing.T) {
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"setup", "/uninstall", "/S"}
	opts := parseArgs()
	if !opts.Uninstall || !opts.Silent {
		t.Fatalf("%+v", opts)
	}
}
