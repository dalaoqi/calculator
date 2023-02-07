package main

import (
	"bufio"
	"calculator/server"
	"calculator/tools"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	envs         map[string]string
	meanServer   server.MeanServer
	medianServer server.MedianServer
	// modeServer   server.ModeServer
	mainServer *http.Server

	router        *mux.Router
	numbersStr    string
	resFromClient string

	input1 chan string
	input2 chan string
	// input3  chan string
	request chan string

	result1 chan string
	result2 chan string
	// result3 chan string

	response chan string
)

func init() {
	loadConfig()
	envs = tools.GetAllEnv()

	meanServer.Init()
	medianServer.Init()
	// modeServer.Init()
	input1, input2 = make(chan string, 1), make(chan string, 1)
	result1, result2 = make(chan string, 1), make(chan string, 1)

	if envs["SERVER_INPUT"] != "keyboard" {
		router = mux.NewRouter()
		router.HandleFunc("/", handler).Methods("POST")
		response, request = make(chan string, 1), make(chan string, 1)
		mainServer = &http.Server{
			Addr:         "127.0.0.1:8081",
			WriteTimeout: time.Second * 300,
			ReadTimeout:  time.Second * 300,
			IdleTimeout:  time.Second * 300,
			Handler:      router,
		}
		go mainServer.ListenAndServe()
	}
}

func main() {
	fmt.Println("Server started...")
	go meanServer.Run(input1, result1)
	go medianServer.Run(input2, result2)
	// go modeServer.Run(input3, result3)
	if <-input1 == "start" && <-input2 == "start" {
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
			// input3 <- numbersStr

			r1 := <-result1
			r2 := <-result2
			// r3 := <-result3
			if r1 == "shutdown\n" && r2 == "shutdown\n" {
				fmt.Println("client1, client2 are shutdown")
				if envs["SERVER_INPUT"] != "keyboard" {
					response <- "shutdown"
				}
				shutdown()
				fmt.Println("Main server is shutdown")
				os.Exit(0)
			}

			resFromClient = fmt.Sprintf("Mean is: %v", r1) + fmt.Sprintf("Median is: %v", r2)

			// if len(r3) > 2 && r3 != "NAN\n" {
			// 	resFromClient += fmt.Sprintf("Mode are: %v", r3)
			// } else {
			// 	resFromClient += fmt.Sprintf("Mode is: %v", r3)
			// }
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

func shutdown() {
	close(input1)
	close(input2)
	// close(input3)

	close(result1)
	close(result2)
	// close(result3)

	if envs["SERVER_INPUT"] != "keyboard" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mainServer.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown: ", err)
		}
		close(request)
		close(response)
	}
}
