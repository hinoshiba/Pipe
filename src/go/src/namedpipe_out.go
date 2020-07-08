package main

import (
	"log"
	"fmt"
	"os"
	"bufio"
	"syscall"
)

var IPC_FIFO_NAME string = "mypipefile"

func main() {
	err := syscall.Mkfifo(IPC_FIFO_NAME, 0666)
	if err != nil {
		log.Println(err)
	}

	f, err := os.OpenFile(IPC_FIFO_NAME, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		l, err := r.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", string(l))
	}
}
