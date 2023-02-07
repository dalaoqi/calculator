package server

import (
	"bufio"
	"fmt"
	"net"
)

var (
	err error
)

type MeanServer struct {
	conn net.Conn
	ln   net.Listener
}

func (s *MeanServer) Init() {
	s.ln, err = net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error starting socket server: " + err.Error())
		panic(err)
	}
}

func (s *MeanServer) Run(numbersStrChan chan string, result chan string) {
	// wait for client
	conn, err := s.ln.Accept()
	if err != nil {
		fmt.Println("Error listening to client: " + err.Error())
	}
	s.conn = conn
	fmt.Println("Client1 is ready...")
	numbersStrChan <- "start"
	for {
		numbersStr := <-numbersStrChan
		go s.receiveData(s.conn, result)
		// fmt.Println("sending to client1... ", numbersStr)
		_, err = fmt.Fprint(s.conn, numbersStr+"\n")
		if err != nil {
			fmt.Println("Client Mean: end sending data")
			return
		}
		if numbersStr == "Q\n" {
			break
		}
	}
}

func (s *MeanServer) receiveData(conn net.Conn, result chan string) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(conn.RemoteAddr().String() + ": client disconnected")
			conn.Close()
			fmt.Println(conn.RemoteAddr().String() + ": end receiving data")
			return
		}
		result <- message
		if message == "shutdown\n" {
			s.shutdown()
			return
		}
	}
}

func (s *MeanServer) shutdown() {
	s.conn.Close()
	s.ln.Close()
	fmt.Println("Mean Service is down")
}
