package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Bind to port 4221
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go requestHandler(conn)
	}
}

func requestHandler(conn net.Conn) {
	defer conn.Close()

	// Read the request
	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		fmt.Println("Error reading request:", err.Error())
		return
	}

	requestLine := string(req[:n])
	fmt.Printf("Request: %s\n", requestLine)

	// Split headers
	splitHeaders := strings.Split(requestLine, "\r\n")
	if len(splitHeaders) < 1 {
		fmt.Println("Malformed request")
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	// Split request line
	splitRequest := strings.Split(splitHeaders[0], " ")
	if len(splitRequest) < 3 {
		fmt.Println("Malformed request line")
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	path := splitRequest[1]

	// Handle different paths
	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(path, "/echo/") {
		parts := strings.SplitN(path, "/", 3)
		if len(parts) < 3 {
			fmt.Println("Malformed echo request")
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			return
		}
		message := parts[2]
		// Construct the response with headers and message
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)
		// Send the response
		conn.Write([]byte(response))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
