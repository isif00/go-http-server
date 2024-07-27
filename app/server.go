package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/config"
	"github.com/codecrafters-io/http-server-starter-go/app/handlers"
)

func main() {
	// Parse command line flags
	flag.IntVar(&config.Port, "port", 4221, "The server port")
	flag.StringVar(&config.Directory, "directory", ".", "The server directory")

	flag.Parse()

	host := fmt.Sprintf("0.0.0.0:%d", config.Port)

	// Bind to port 4221
	l, err := net.Listen("tcp", host)
	fmt.Println("Listening on " + host)

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handlers.RequestHandler(conn)
	}
}
