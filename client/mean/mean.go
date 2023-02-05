package main

import (
	"bufio"
	"calculator/tools"
	"fmt"
	"net"
	"strconv"
	"time"
)

var (
	connected bool
	result    string
	numbers   []int
)

func main() {
	fmt.Println("Mean Client started...")
	for {
		if !connected {
			conn, err := net.Dial("tcp", "127.0.0.1:8000")
			if err != nil {
				fmt.Println(err.Error())
				time.Sleep(time.Duration(5) * time.Second)
				continue
			}
			fmt.Println(conn.RemoteAddr().String() + ": connected")
			connected = true
			go receiveData(conn)
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func receiveData(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Server: disconnected")
			conn.Close()
			connected = false
			fmt.Println("Server:: end receiving data")
			return
		}
		fmt.Print(conn.RemoteAddr().String() + ": received " + message)

		numbers = tools.Str2slice(message)
		result = mean(numbers)
		_, err = fmt.Fprint(conn, result)
		if err != nil {
			fmt.Println("Server: end sending data")
			return
		}
	}
}

func mean(numbers []int) string {
	if len(numbers) == 0 {
		return "NAN\n"
	}
	var sum float64
	var res string
	for _, i := range numbers {
		sum += float64(i)
	}
	if sum == 0 {
		res = "0\n"
	} else {
		res = strconv.FormatFloat(sum/float64(len(numbers)), 'f', 13, 64) + "\n"
	}
	return res
}
