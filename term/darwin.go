//go:build darwin

package term

import (
	"fmt"
	"syscall"
	"unsafe"
)

const ResetColourSequence = "\033[0m" // anything after this will have no custom colouring applied

func EnableRawMode(fd int) (*syscall.Termios, error) {
	var termios syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCGETA, uintptr(unsafe.Pointer(&termios)), 0, 0, 0); err != 0 {
		return nil, err
	}

	origState := termios

	newState := termios
	newState.Lflag &^= syscall.ECHO | syscall.ICANON

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCSETA, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		return nil, err
	}

	return &origState, nil
}

func Restore(fd int, origState *syscall.Termios) error {
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCSETA, uintptr(unsafe.Pointer(origState)), 0, 0, 0); err != 0 {
		return err
	}
	return nil
}

func Clear() {
	fmt.Print("\033[H\033[2J\033[3J")
}
