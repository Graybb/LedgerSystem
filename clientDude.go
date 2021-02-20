package peerToPeer

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func AskForInput(text string) string {
	fmt.Printf("%s: ", text)
	return ReadInput()
}

func ReadInput() string {
	response, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Trim(response, "\n")
}

func Connect(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to connect")
		return
	}

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

func main() {
	ip := AskForInput("Enter IP address")
	port := AskForInput("Enter port")
	address := ip + ":" + port
	//address := "127.0.0.1:8080"

	fmt.Printf("Connecting to %s\n", address)

	Connect(address)
}