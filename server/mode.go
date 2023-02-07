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
	SEND_KEY    = 1000
	RECEIVE_KEY = 1001
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
	sendId      uintptr
	sendAddr    uintptr
	receiveId   uintptr
	receiveAddr uintptr
}

func (s *ModeServer) Init() {
	s.sendId, s.sendAddr = createMemorySegment(SEND_KEY)
	s.receiveId, s.receiveAddr = createMemorySegment(RECEIVE_KEY)
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
		fmt.Println("sending to modeClient... ", numbersStr)
		if numbersStr == "Q" {
			fmt.Printf("break\n")
			break
		}
	}
}

func createMemorySegment(key int) (shmId, addr uintptr) {
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
		deleteMemorySegment(shmId, addr)
	}
	if err != nil {
		deleteMemorySegment(shmId, addr)
	}
	return shmId, addr
}

func deleteMemorySegment(id, addr uintptr) (int, error) {
	result, _, errno := syscall.Syscall(syscall.SYS_SHMDT, addr, 0, 0)
	if int(result) == -1 {
		return -1, errno
	}

	result, _, errno = syscall.Syscall(syscall.SYS_SHMCTL, id, 0, 0)
	if int(result) == -1 {
		return -1, errno
	}
	return int(result), nil
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
		var receiveData = data{s.receiveAddr, int(SEGSIZE), int(SEGSIZE)}
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
		var receiveData = data{s.receiveAddr, int(SEGSIZE), int(SEGSIZE)}
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
		fmt.Printf("server read: %v", sharedText.Written)
		if sharedText.Written {
			fmt.Printf("server write: %v", sharedText.Text)
			result <- sharedText.Text
			sharedText.Written = false
			err := enc.Encode(sharedText)
			if err != nil {
				log.Fatal("encode error:", err)
			}
			copy(shmData, buf.Bytes())
			if sharedText.Text == "shutdown\n" {
				s.shutdown()
				break
			}
		}
	}
}

func (s *ModeServer) shutdown() {
	_, err := deleteMemorySegment(s.sendId, s.sendAddr)
	if err != nil {
		log.Fatal("deleteMemorySegment error:", err)
	}
	_, err = deleteMemorySegment(s.receiveId, s.receiveAddr)
	if err != nil {
		log.Fatal("deleteMemorySegment error:", err)
	}
	fmt.Println("client3 is shutdown")
}
