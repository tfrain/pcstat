package main

import (
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// adapted from https://groups.google.com/d/msg/golang-nuts/8d4pOPmSL9Q/H6WUqbGNELEJ
type winsize struct {
	ws_row, ws_col       uint16
	ws_xpixel, ws_ypixel uint16
}

func getwinsize() winsize {
	ws := winsize{}
	_, _, err := unix.Syscall(syscall.SYS_IOCTL,
		uintptr(0), uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))
	if err != 0 {
		log.Fatalf("TIOCGWINSZ failed to get terminal size: %s\n", err)
	}
	return ws
}
