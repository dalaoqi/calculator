package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/gen2brain/shm"
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

const (
	SEGSIZE     = 65536
	FLAG        = 00001000 | 0777
	SEND_KEY    = 100
	RECEIVE_KEY = 101
)

var (
	sendAddr    uintptr
	recevieAddr uintptr
	result      chan string
	modeStr     string
)

func init() {
	sendAddr = createMemorySegment(SEND_KEY)
	recevieAddr = createMemorySegment(RECEIVE_KEY)
	result = make(chan string)
}

func main() {
	fmt.Println("Mode Client started...")
	go writer(result)
	for {
		var receiveData = data{sendAddr, int(SEGSIZE), int(SEGSIZE)}
		shmData := *(*[]byte)(unsafe.Pointer(&receiveData))
		dataFromServer := make([]byte, len(shmData))
		copy(dataFromServer, shmData)

		var sharedText sharedStat
		var buf bytes.Buffer
		buf.Write(dataFromServer)
		enc := gob.NewEncoder(&buf)
		dec := gob.NewDecoder(&buf)
		err := dec.Decode(&sharedText)
		if err != nil {
			continue
		}
		if sharedText.Written {
			if sharedText.Text == "hello\n" {
				result <- "HELLO\n"
			} else {
				fmt.Printf("receive %v from server\n", sharedText.Text)
				numbersStr := strings.Fields(sharedText.Text)
				var numbers []int
				for _, numberStr := range numbersStr {
					number, err := strconv.Atoi(numberStr)
					if err != nil {
						fmt.Println("Invalid input. Please enter only numbers.")
						continue
					}
					numbers = append(numbers, number)
				}
				modeStr = mode(numbers)
				fmt.Printf("complete: %v", modeStr)
				result <- modeStr
			}
			sharedText.Written = false
			err := enc.Encode(sharedText)
			if err != nil {
				log.Fatal("encode error:", err)
			}
			copy(shmData, buf.Bytes())
		}
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

func mode(numbers []int) string {
	if len(numbers) == 0 {
		return "NAN\n"
	}
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	count := 0
	countMap := make(map[int]int)
	resInt := make([]int, 0)
	resStr := make([]string, 0)
	for _, number := range numbers {
		countMap[number] += 1
		count = max(count, countMap[number])
	}

	for k, v := range countMap {
		if v == count {
			resInt = append(resInt, k)
		}
	}

	sort.Ints(resInt)

	for _, v := range resInt {
		resStr = append(resStr, strconv.Itoa(v))
	}

	return strings.Join(resStr, " ") + "\n"
}

func writer(result chan string) {
	for {
		numbersStr := <-result
		h := data{recevieAddr, int(SEGSIZE), int(SEGSIZE)}
		var buf bytes.Buffer
		sharedText := sharedStat{true, numbersStr}
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(sharedText)
		if err != nil {
			log.Fatal("encode error:", err)
		}
		shmData := *(*[]byte)(unsafe.Pointer(&h))
		copy(shmData, buf.Bytes())
		time.Sleep(1 * time.Second)
	}
}
