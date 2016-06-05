// +build linux

package main

import (
	"os"
	"unsafe"

	"syscall"
)

func TermWidth() (width int, ok bool) {
	var winsize struct {
		wsRow    uint16
		wsCol    uint16
		wsXpixel uint16
		wsYpixel uint16
	}

	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&winsize))); errno != 0 {
		return 0, false
	}
	return int(winsize.wsCol), true
}
