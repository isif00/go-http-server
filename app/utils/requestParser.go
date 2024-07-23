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

	// Split user agent
	var userAgent string
	for _, header := range splitHeaders {
		if strings.HasPrefix(header, "User-Agent:") {
			userAgentParts := strings.SplitN(header, " ", 2)
			if len(userAgentParts) == 2 {
				userAgent = userAgentParts[1]
			}
		}
	}

	req.UserAgent = userAgent
	req.Path = splitRequest[1]

	return &req, nil
}
