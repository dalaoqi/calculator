package main

import (
	"bufio"
	"calculator/server"
	"fmt"
	"os"
	"strings"
)

var (
	meanServer   server.MeanServer
	medianServer server.MedianServer
	modeServer   server.ModeServer

	input1 chan string
	input2 chan string
	input3 chan string

	result1 chan string
	result2 chan string
	result3 chan string
)

func init() {
	meanServer.Init()
	medianServer.Init()
	modeServer.Init()
	input1, input2, input3 = make(chan string), make(chan string), make(chan string)
	result1, result2, result3 = make(chan string), make(chan string), make(chan string)
}

func main() {
	fmt.Println("Server started...")
	go meanServer.Run(input1, result1)
	go medianServer.Run(input2, result2)
	go modeServer.Run(input3, result3)
	if <-input1 == "start" && <-input2 == "start" && <-input3 == "start" {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("You can type intergers and then click [ENTER].  Clients will show the mean, median, and mode of the input values.\n")

			input, _ := reader.ReadString('\n')
			numbersStr := strings.Join(strings.Fields(input), " ")

			input1 <- numbersStr
			input2 <- numbersStr
			input3 <- numbersStr

			fmt.Printf("Mean is: %v", <-result1)
			fmt.Printf("Median is: %v", <-result2)
			r3 := <-result3
			if len(r3) > 2 && r3 != "NAN\n" {
				fmt.Printf("Mode are: %v", r3)
			} else {
				fmt.Printf("Mode is: %v", r3)
			}
		}
	}
}
