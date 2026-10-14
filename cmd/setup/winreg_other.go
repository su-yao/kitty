//go:build !windows

package main

import "errors"

func createShortcut(dir, name, target, workDir, desc string) error {
	return errors.New("shortcuts require Windows")
}

func removeShortcut(dir, name string) {}

func writeUninstallKey(dest, uninst string) error {
	return errors.New("registry requires Windows")
}

func deleteUninstallKey() error { return nil }

func readInstallDir() (string, error) {
	return "", errors.New("registry requires Windows")
}

func selfDelete(dest string) {}

func launchAdmin(exe string) error {
	return errors.New("launch requires Windows")
}
