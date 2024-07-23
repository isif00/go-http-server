package utils

import (
	"net"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/types"
)

func ParseRequest(requestLine string, conn net.Conn) (*types.HttpRequest, error) {
	var req types.HttpRequest

	// Split headers
	splitHeaders := strings.Split(requestLine, "\r\n")
	if len(splitHeaders) < 1 {
		conn.Write([]byte(GetStatus(400, "Bad Request")))
	}

	// Split request line
	splitRequest := strings.Split(splitHeaders[0], " ")
	if len(splitRequest) < 3 {
		conn.Write([]byte(GetStatus(400, "Bad Request")))
	}

	req.Path = splitRequest[1]

	return &req, nil
}
