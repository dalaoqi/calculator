package server

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
)

const (
	SEND_PIPEFILE    = "/tmp/sendPipe.log"
	RECEIVE_PIPELINE = "/tmp/receivePipe.log"
)

type MedianServer struct {
}

func (s *MedianServer) Init() {
	os.Remove(SEND_PIPEFILE)
	err = syscall.Mkfifo(SEND_PIPEFILE, 0666)
	if err != nil {
		log.Fatalf("Make %v error:%v", SEND_PIPEFILE, err)
	}

	os.Remove(RECEIVE_PIPELINE)
	err = syscall.Mkfifo(RECEIVE_PIPELINE, 0666)
	if err != nil {
		log.Fatalf("Make %v error:%v", RECEIVE_PIPELINE, err)
	}
}

func (s *MedianServer) Run(numbersStrChan chan string, result chan string) {
	path := SEND_PIPEFILE

	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		log.Panicln(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	if s.test2() {
		fmt.Println("Client2 is ready...")
		numbersStrChan <- "start"
	}
	go s.receivePipe(result)
	for {
		numbersStr := <-numbersStrChan
		// fmt.Println("sending to client3... ", numbersStr)
		if _, err := f.WriteString(numbersStr + "\n"); err != nil {
			log.Panicln(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *MedianServer) test2() bool {
	file, err := os.OpenFile(RECEIVE_PIPELINE, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		log.Fatalf("Open %v error: %v", RECEIVE_PIPELINE, err)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		if line == "hello\n" {
			return true
		}
	}
}

func (s *MedianServer) receivePipe(result chan string) {
	file, err := os.OpenFile(RECEIVE_PIPELINE, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		log.Fatalf("Open %v error: %v", RECEIVE_PIPELINE, err)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		result <- line
	}
}
