package peerToPeer

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
			fmt.Println("From " + otherEnd  + " to " + myEnd + ": " + string(msg))
			channel <- msg
		}
	}
}



func Connect(address string, chanstrings chan string, ) { //tries to connect to an node
	conn, err := net.Dial("tcp", address)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to connect")
		return
	}

	fmt.Println("Connected")

	//puts the main connection in the connection map. and sets its value to true
	connectionMap[conn] = true

	handleConnection(conn, chanstrings)


}

func openForTCPConnection(channel chan string) { //accepts to all
			ln, _ := net.Listen("tcp", ":")  //listens on a random port
	    	defer ln.Close()
			fmt.Println("Listening for connection... on port"+ln.Addr().String())
	    	for  {
	    		conn, _ := ln.Accept()
	    		connectionMap[conn] = true
	    		fmt.Println("Got a connection from address:  " + conn.RemoteAddr().String())
	    		go handleConnection(conn, channel)
	    	}
}

func main() {
	quit = true
	//prints own address and port
	name, _ := os.Hostname()
	myAddrs, _ := net.LookupHost(name)
	fmt.Println("Name: " + name)
	fmt.Println(myAddrs)

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
	go handleChannelOutPut(c)

	fmt.Printf("Connecting to %s\n", address)

	go Connect(address,c)

	go openForTCPConnection(c)
	askForMsg(c)

}

func askForMsg(chanstrings chan string) {
	for {
		fmt.Print("> ")
		input := ReadInput()
		if input == "quit" { quit = false
			return }
		//sends to the server and hopefully each child
		chanstrings <- input

	}
}



func handleChannelOutPut(c chan string) {
	for{
		msg := <- c
		for key,_ := range connectionMap {
				key.Write([]byte(msg))
				fmt.Println(connectionMap)
		}

	}
}


