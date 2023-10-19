package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"

	"github.com/spf13/pflag"
)

func main() {
	file := pflag.StringP("file", "f", "-", "path to file to read/write to, or '-' for stdin")
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"run '%s get' to get current layout\n"+
				"run '%s set' to read and set layout from stdin\n"+
				"use '-f [path]' to use a file instead of stdin\n\n"+
				"find all layout codes here:\n"+
				"https://learn.microsoft.com/en-us/windows-hardware/manufacture/desktop/windows-language-pack-default-values",
			os.Args[0], os.Args[0])
		os.Exit(1)
	}
	pflag.Parse()

	switch pflag.Arg(0) {
	case "get":
		w := os.Stdout
		var err error
		if *file != "-" {
			w, err = os.Create(*file)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open %q for writing: %v", *file, err)
			os.Exit(1)
		}

		layout, err := GetKeyboardLayout()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get keyboard layout: %v", err)
			os.Exit(1)
		}

		fmt.Fprint(w, layout)
		os.Exit(0)

	case "set":
		r := os.Stdin
		var err error
		if *file != "-" {
			r, err = os.Open(*file)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open %q for reading: %v", *file, err)
			os.Exit(1)
		}

		var buf [8]byte
		_, err = io.ReadFull(r, buf[:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read layout: %v", err)
			os.Exit(1)
		}

		err = SetKeyboardLayout(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to set keyboard layout: %v", err)
			os.Exit(1)
		}

		os.Exit(0)

	default:
		pflag.Usage()
	}
}

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	getKeyboardLayout  = user32.NewProc("GetKeyboardLayoutNameA")
	loadKeyboardLayout = user32.NewProc("LoadKeyboardLayoutA")
	sendMessage        = user32.NewProc("SendMessageA")
)

func GetKeyboardLayout() (string, error) {
	var buf [8]byte

	r, _, err := getKeyboardLayout.Call(
		uintptr(unsafe.Pointer(&buf[0])),
	)
	if r == 0 {
		return "00000000", err
	}

	return string(buf[:]), nil
}

func SetKeyboardLayout(layout [8]byte) error {
	locale, _, err := loadKeyboardLayout.Call(
		uintptr(unsafe.Pointer(&layout[0])),
		uintptr(1),
	)
	if locale == 0 {
		return err
	}

	if r, _, err := sendMessage.Call(
		0xffff, // Broadcast
		0x0050, // change request
		0,      // ??
		locale,
	); r == 0 {
		return err
	}

	return nil
}
