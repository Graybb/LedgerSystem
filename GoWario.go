
package peerToPeer
import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)



func handleConnection(conn net.Conn, c chan string) {
	defer conn.Close()
	myEnd := conn.LocalAddr().String()
	otherEnd := conn.RemoteAddr().String()
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if (err != nil) {
			fmt.Println("Ending session with " + otherEnd)
			return
		} else {
			fmt.Print("From " + otherEnd  + " to " + myEnd + ": " + string(msg))
			titlemsg := conn.RemoteAddr().String() + "18081" + strings.Title(msg)
			conn.Write([]byte(titlemsg))
			c <- titlemsg

		}
	}
}


func main() {
	name, _ := os.Hostname()
	addrs, _ := net.LookupHost(name)
	fmt.Println("Name: " + name)

	c := make (chan string)
	for indx, addr := range addrs {
		fmt.Println("Address number " + strconv.Itoa(indx) + ": " + addr)
	}  // writes the local adress of server



	ln, _ := net.Listen("tcp", ":")
	defer ln.Close()
	for {
		fmt.Println("Listening for connection... on port"+ln.Addr().String())
		conn, _ := ln.Accept()
		fmt.Println("Got a connection...")
		go handleConnection(conn, c)

	}

}



