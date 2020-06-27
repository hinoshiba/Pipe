package main

import (
	"io"
	"net"
	"os"
)

const (
	SOCK_FILE = "./mysocket.sock"
)

func socket() error {
	stdin := os.Stdin
	stdout := os.Stdout

	b_ch := make(chan []byte)
	err_ch := make(chan error)

	sv, err := net.Listen("unix", SOCK_FILE)
	if err != nil {
		return err
	}
	defer os.Remove(SOCK_FILE)

	go func() {
		for {
			sv_session, err := sv.Accept()
			if err != nil {
				err_ch <- err
				return
			}

			go func() {
				defer sv_session.Close()

				for {
					buf := make([]byte, 8)
					size, err := sv_session.Read(buf)
					if err != nil {
						if err != io.EOF {
							err_ch <- err
							return
						}

						break
					}

					stdout.Write(buf[:size])
				}
			}()
		}
	}()

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
			session, err := net.Dial("unix", SOCK_FILE)
			if err != nil {
				err_ch <- err
				return
			}
			defer session.Close()

			_, err = session.Write(bytes)
			if err != nil {
				err_ch <- err
				return
			}
			if err = session.(*net.UnixConn).CloseWrite(); err != nil {
				err_ch <- err
				return
			}
		}
	}()

	for {
		err := <- err_ch
		return err
	}
}

func main() {
	if err := socket(); err != nil {
		panic(err)
	}
}
