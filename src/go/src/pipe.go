package main

import (
	"os"
)

func pipe() error {
	stdin := os.Stdin
	stdout := os.Stdout

	b_ch := make(chan []byte)
	err_ch := make(chan error)

	go func() {
		for {
			buf := make([]byte, 8, 8)
			n, err := stdin.Read(buf)
			if err != nil {
				err_ch <- err
				return
			}
			if n < 1 {
				continue
			}
			b_ch <- buf
		}
	}()

	go func() {
		for {
			bytes := <- b_ch
			if bytes == nil {
				continue
			}

			stdout.Write(bytes)
		}
	}()

	for {
		err := <- err_ch
		return err
	}
}

func main() {
	if err := pipe(); err != nil {
		panic(err)
	}
}
