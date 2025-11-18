//go:build linux

package term

import (
	"fmt"
	"syscall"
	"unsafe"
)

func EnableRawMode(fd int) (*syscall.Termios, error) {
	var termios syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TCGETS, uintptr(unsafe.Pointer(&termios)), 0, 0, 0); err != 0 {
		return nil, err
	}

	origState := termios

	newState := termios
	newState.Lflag &^= syscall.ECHO | syscall.ICANON

	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&newState)), 0, 0, 0)
	if err != 0 {
		return nil, err
	}
	return &origState, nil
}

func Restore(fd int, origState *syscall.Termios) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TCSETS, uintptr(unsafe.Pointer(&origState)), 0, 0, 0)
	return err
}

func Clear() {
	fmt.Print("\033[H\033[2J\033[3J")
}
