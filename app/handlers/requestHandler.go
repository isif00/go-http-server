package handlers

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

func RequestHandler(conn net.Conn) {
	// Read the request
	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		fmt.Println("Error reading request:", err.Error())
		return
	}

	requestLine := string(req[:n])
	fmt.Printf("Request: %s\n", requestLine)

	parsedRequest, nil := utils.ParseRequest(requestLine, conn)
	if nil != nil {
		conn.Write([]byte(utils.GetStatus(400, "Bad Request")))
		return
	}

	var response string

	switch path := parsedRequest.Path; {

	case path == "/":
		response = utils.GetStatus(200, "OK\r\n")

	case path == "/echo":
		parts := strings.SplitN(path, "/", 3)
		if len(parts) < 3 {
			conn.Write([]byte(utils.GetStatus(400, "Bad Request\r\n")))
			return
		}
		message := parts[2]
		response = fmt.Sprintf("%sContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", utils.GetStatus(200, "OK"), len(message), message)

	case path == "/user-agent":
		response = fmt.Sprintf("%sContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", utils.GetStatus(200, "OK"), len(parsedRequest.UserAgent), parsedRequest.UserAgent)

	default:
		response = utils.GetStatus(404, "Not Found\r\n")
	}

	conn.Write([]byte(response))
	conn.Close()
}
