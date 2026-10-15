//go:build !windows

package main

import "log"

func notify(title, msg string) {
	log.Println(title + ": " + msg)
}
