package main

import (
	"fmt"
	"net"
)

func main() {
	// Connect to a TCP server
	conn, err := net.Dial("tcp", "localhost:6969")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	// Send and receive data
	fmt.Fprintf(conn, "GET / HTTP/1.1 \nHost: example.com\n\n")
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println(string(buf[:n]))
}
