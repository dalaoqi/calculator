package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"github.com/gen2brain/shm"
)

const (
	SEGSIZE     = 65536
	FLAG        = 00001000 | 0777
	SEND_KEY    = 100
	RECEIVE_KEY = 101
)

type sharedStat struct {
	Written bool
	Text    string
}

type data struct {
	addr uintptr
	len  int
	cap  int
}
type ModeServer struct {
	sendAddr    uintptr
	recevieAddr uintptr
}

func (s *ModeServer) Init() {
	s.sendAddr = createMemorySegment(SEND_KEY)
	s.recevieAddr = createMemorySegment(RECEIVE_KEY)
}

func (s *ModeServer) Run(numbersStrChan chan string, result chan string) {
	if s.test3() {
		fmt.Println("Client3 is ready...")
		numbersStrChan <- "start"
	}
	go s.receiveSHM(result)
	for {
		numbersStr := <-numbersStrChan
		var buf bytes.Buffer
		sharedText := sharedStat{true, numbersStr}
		enc := gob.NewEncoder(&buf)
		err = enc.Encode(sharedText)
		if err != nil {
			log.Fatal("encode error:", err)
		}

		h := data{s.sendAddr, int(SEGSIZE), int(SEGSIZE)}

		shmData := *(*[]byte)(unsafe.Pointer(&h))
		copy(shmData, buf.Bytes())
		// fmt.Println("sending to modeClient... ", numbersStr)
	}
}

func createMemorySegment(key int) (addr uintptr) {
	shmId, _, errno := syscall.Syscall(syscall.SYS_SHMGET, uintptr(int32(key)), uintptr(int32(SEGSIZE)), uintptr(int32(FLAG)))
	// fmt.Printf("id1: %v\n", shmId)
	if int(shmId) == -1 {
		panic(errno)
	}
	addr, _, errno = syscall.Syscall(syscall.SYS_SHMAT, uintptr(int32(shmId)), uintptr(int32(0)), uintptr(int32(0)))
	// fmt.Printf("addr1: %v\n", addr)
	if int(addr) == -1 {
		panic(errno)
	}
	length, err := shm.Size(int(shmId))
	// fmt.Printf("size of memory segment: %v\n", length)
	if length != SEGSIZE {
		fmt.Printf("SIZE Error\n")
		syscall.Syscall(syscall.SYS_SHMDT, addr, 0, 0)
	}
	if err != nil {
		syscall.Syscall(syscall.SYS_SHMDT, addr, 0, 0)
	}
	return addr
}

func (s *ModeServer) test3() bool {
	var buf bytes.Buffer
	sharedText := sharedStat{true, "hello\n"}
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(sharedText)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	h := data{s.sendAddr, int(SEGSIZE), int(SEGSIZE)}

	shmData := *(*[]byte)(unsafe.Pointer(&h))
	copy(shmData, buf.Bytes())

	for {
		var receiveData = data{s.recevieAddr, int(SEGSIZE), int(SEGSIZE)}
		shmData := *(*[]byte)(unsafe.Pointer(&receiveData))
		dataFromClient := make([]byte, len(shmData))
		copy(dataFromClient, shmData)

		var sharedText sharedStat
		var buf bytes.Buffer
		buf.Write(dataFromClient)
		enc := gob.NewEncoder(&buf)
		dec := gob.NewDecoder(&buf)
		err := dec.Decode(&sharedText)
		if err != nil {
			continue
		}
		if sharedText.Written {
			sharedText.Written = false
			err := enc.Encode(sharedText)
			if err != nil {
				log.Fatal("encode error:", err)
			}
			copy(shmData, buf.Bytes())
			if sharedText.Text == "HELLO\n" {
				return true
			}
		}
	}
}

func (s *ModeServer) receiveSHM(result chan string) {
	for {
		var receiveData = data{s.recevieAddr, int(SEGSIZE), int(SEGSIZE)}
		shmData := *(*[]byte)(unsafe.Pointer(&receiveData))
		dataFromClient := make([]byte, len(shmData))
		copy(dataFromClient, shmData)

		var sharedText sharedStat
		var buf bytes.Buffer
		buf.Write(dataFromClient)
		enc := gob.NewEncoder(&buf)
		dec := gob.NewDecoder(&buf)
		err := dec.Decode(&sharedText)
		if err != nil {
			continue
		}
		if sharedText.Written {
			result <- sharedText.Text
			sharedText.Written = false
			err := enc.Encode(sharedText)
			if err != nil {
				log.Fatal("encode error:", err)
			}
			copy(shmData, buf.Bytes())
		}
	}
}
