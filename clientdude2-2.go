package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)
var quit bool
 var connectionMap map[net.Conn] bool

func AskForInput(text string) string {
	fmt.Printf("%s: ", text)
	return ReadInput()
}

func ReadInput() string {
	response, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Trim(response, "\n")
}


func handleConnection(conn net.Conn, channel chan string) { //reciveing connections for peer to peer
	defer conn.Close()
	defer delete(connectionMap,conn)
	myEnd := conn.LocalAddr().String()
	otherEnd := conn.RemoteAddr().String()
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if (err != nil) {
			fmt.Println("Ending session with " + otherEnd)
			return
		} else {
			fmt.Print("From " + otherEnd  + " to " + myEnd + ": " + string(msg))
			titlemsg := strings.Title(msg)
			channel <- titlemsg

		}
	}
}



func Connect(address string) { //tries to connect to an node
	conn, err := net.Dial("tcp", address)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to connect")
		return
	}
	handleConnection()

	fmt.Println("Connected")

	defer conn.Close()
	for {
		fmt.Print("> ")
		input := ReadInput()
		if input == "quit" { return }
		_, _ = fmt.Fprintln(conn, input)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil { return }
		fmt.Println("From server: " + message)
	}
}

func openForTCPConnection(channel chan string){       //accepts to all
	    ln, _ := net.Listen("tcp", ":")  //listens on a random port
	    	defer ln.Close()
	    	for  {
	    		fmt.Println("Listening for connection... on port"+ln.Addr().String())
	    		conn, _ := ln.Accept()
	    		connectionMap[conn] = true
	    		fmt.Println("Got a connection from address:  " + conn.RemoteAddr().String())
	    		go handleConnection(conn, channel)
	    	}
}

func main() {
	//prints own address and port
	name, _ := os.Hostname()
	MyAddrs, _ := net.LookupHost(name)
	fmt.Println("Name: " + name)
	fmt.Println(MyAddrs)

	//for indx, MyAddrs := range MyAddrs {
	//	fmt.Println("Address number " + strconv.Itoa(indx) + ": " + MyAddrs)
	//}
	connectionMap = make(map[net.Conn] bool)

	//ask for another port and address
	ip := AskForInput("Enter IP address")
	port := AskForInput("Enter port")
	address := ip + ":" + port

	//address := "127.0.0.1:8080"
	c := make(chan string)
	fmt.Printf("Connecting to %s\n", address)
	go Connect(address)
	go handleChannelOutPut(c)

	openForTCPConnection(c)


}

func handleChannelOutPut(c chan string) {
	for{
		msg := <- c
		for key,value := range connectionMap {
			if value {
				key.Write([]byte(msg))
			}
		}

	}


}