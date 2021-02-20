package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var connectionMap map[net.Conn]bool
var messagesMap map[string]bool

func AskForInput(text string) string {
	fmt.Printf("%s: ", text)
	return ReadInput()
}

func ReadInput() string {
	response, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Trim(response, "\n")

}

func handleConnection(conn net.Conn, channel chan string) { //reciveing connections for peer to peer
	//string fluff
	//myEnd := conn.LocalAddr().String()
	otherEnd := conn.RemoteAddr().String()

	msg, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Println("Ending session with " + otherEnd)
		return
	} else {
		go handleConnection(conn, channel)
		_, ok := messagesMap[msg]
		if !ok {
			fmt.Println("From : " + string(msg))
			channel <- msg
		}

	}

}

func Connect(address string, chanstrings chan string) { //tries to connect to an node
	conn, err := net.Dial("tcp", address)

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to connect")
		for {
			fmt.Print("> ")
			input := ReadInput()
			if input == "quit" {
				return
			}
			chanstrings <- input
		}
		return
	}
	fmt.Println("Connected")
	connectionMap[conn] = true
	go handleConnection(conn, chanstrings)

	for {
		fmt.Print("> ")
		input := ReadInput()
		if input == "quit" {
			return
		}

		chanstrings <- input
	}

}

func openForTCPConnection(channel chan string) { //accepts to all
	ln, _ := net.Listen("tcp", ":") //listens on a random port
	defer ln.Close()
	for {
		fmt.Println("Listening for connection on port: " + ln.Addr().String())
		conn, _ := ln.Accept()
		connectionMap[conn] = true
		fmt.Println("Got a connection from address:  " + conn.RemoteAddr().String())
		go handleConnection(conn, channel)
	}
}

func main() {
	//prints own address and port
	name, _ := os.Hostname()
	myAddrs, _ := net.LookupHost(name)
	fmt.Println("Name: " + name)
	fmt.Println(myAddrs)

	messagesMap = make(map[string]bool)
	connectionMap = make(map[net.Conn]bool)

	//ask for another port and address
	ip := AskForInput("Enter IP address")
	port := AskForInput("Enter port")
	address := ip + ":" + port

	//address := "127.0.0.1:8080"
	c := make(chan string)

	go handleChannelOutPut(c)
	go Connect(address, c)
	openForTCPConnection(c)

}

func handleChannelOutPut(c chan string) {
	for {
		msg := <-c
		titlemsg := strings.Title(msg)
		_, ok := messagesMap[titlemsg]
		if !ok {
			messagesMap[titlemsg] = true
			for key, _ := range connectionMap {
				key.Write([]byte(titlemsg + "\n"))
				fmt.Fprintln(key, titlemsg)
			}
		}
	}
}
