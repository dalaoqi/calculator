package main

import (
	"bufio"
	"calculator/server"
	"calculator/tools"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	envs         map[string]string
	meanServer   server.MeanServer
	medianServer server.MedianServer
	modeServer   server.ModeServer

	numbersStr    string
	resFromClient string

	input1  chan string
	input2  chan string
	input3  chan string
	request chan string

	result1 chan string
	result2 chan string
	result3 chan string

	response chan string
)

func init() {
	loadConfig()
	envs = tools.GetAllEnv()

	meanServer.Init()
	medianServer.Init()
	modeServer.Init()
	input1, input2, input3, request = make(chan string), make(chan string), make(chan string), make(chan string)
	result1, result2, result3 = make(chan string), make(chan string), make(chan string)

	if envs["SERVER_INPUT"] != "keyboard" {
		response = make(chan string)
		http.HandleFunc("/", handler)
		go http.ListenAndServe(":8081", nil)
	}
}

func main() {
	fmt.Println("Server started...")
	go meanServer.Run(input1, result1)
	go medianServer.Run(input2, result2)
	go modeServer.Run(input3, result3)
	if <-input1 == "start" && <-input2 == "start" && <-input3 == "start" {
		for {
			if envs["SERVER_INPUT"] == "keyboard" {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("You can type intergers and then click [ENTER].  Clients will show the mean, median, and mode of the input values.\n")

				input, _ := reader.ReadString('\n')
				numbersStr = strings.Join(strings.Fields(input), " ")
			} else {
				fmt.Print("Waiting for the request...\n")
				numbersStr = <-request
				fmt.Printf("Get a sequence [%v]\n", numbersStr)
			}
			input1 <- numbersStr
			input2 <- numbersStr
			input3 <- numbersStr

			resFromClient = fmt.Sprintf("Mean is: %v", <-result1) + fmt.Sprintf("Median is: %v", <-result2)

			r3 := <-result3
			if len(r3) > 2 && r3 != "NAN\n" {
				resFromClient += fmt.Sprintf("Mode are: %v", r3)
			} else {
				resFromClient += fmt.Sprintf("Mode is: %v", r3)
			}
			fmt.Println("-------------------")
			fmt.Print(resFromClient)
			fmt.Println("-------------------")

			if envs["SERVER_INPUT"] != "keyboard" {
				response <- resFromClient
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	str, _ := io.ReadAll(r.Body)
	request <- string(str)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(<-response))
}

func loadConfig() {
	fmt.Printf("env:\n")
	if len(os.Args) > 1 {
		f, _ := filepath.Abs(os.Args[1])
		err := tools.LoadConfigFromFile(f)
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		f, _ := filepath.Abs("./conf/env.conf")
		err := tools.LoadConfigFromFile(f)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}
