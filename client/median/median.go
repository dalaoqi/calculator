package main

import (
	"bufio"
	"calculator/tools"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"syscall"
	"time"
)

const (
	SEND_PIPEFILE    = "/tmp/sendPipe.log"
	RECEIVE_PIPELINE = "/tmp/receivePipe.log"
)

var (
	result    chan string
	medianStr string
	numbers   []int
)

func init() {
	result = make(chan string)
}

func main() {
	path := SEND_PIPEFILE

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		// try to create a named pipe
		if err := syscall.Mkfifo(path, 0600); err != nil {
			log.Panicln(err)
		}
	}

	f, err := os.Open(path)
	if err != nil {
		log.Panicln(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	fmt.Println("Median Client started...")
	go writer(result)
	result <- "hello\n"
	for {
		message, err := bufio.NewReader(f).ReadString('\n')
		if err != nil {
			log.Panicln(err)
		}

		numbers = tools.Str2slice(message)
		medianStr = median(numbers)
		fmt.Printf("Median is: %v", medianStr)
		result <- medianStr
		time.Sleep(1 * time.Second)
	}
}

func writer(result chan string) {
	path := RECEIVE_PIPELINE

	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		log.Panicln(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	for {
		numStr := <-result
		if _, err := f.WriteString(numStr); err != nil {
			log.Panicln(err)
		}
	}
}

func median(numbers []int) string {
	sort.Ints(numbers)
	n := len(numbers)
	var res string
	if n == 0 {
		return "NAN\n"
	} else if n%2 == 0 {
		p := n / 2
		res = strconv.FormatFloat((float64(numbers[p-1])+float64(numbers[p]))/2, 'f', 13, 64)
	} else {
		res = strconv.Itoa(numbers[n/2])
	}

	return res + "\n"
}
