package utils

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/types"
)

func ParseRequest(requestLine string, conn net.Conn) (*types.HttpRequest, error) {
	var req types.HttpRequest

	// Split headers and body
	parts := strings.SplitN(requestLine, "\r\n\r\n", 2)
	headersPart := parts[0]
	var bodyPart string
	if len(parts) > 1 {
		bodyPart = parts[1]
	}

	// Split headers
	splitHeaders := strings.Split(headersPart, "\r\n")
	if len(splitHeaders) < 1 {
		conn.Write([]byte(GetStatus(400, "Bad Request")))
		return nil, fmt.Errorf("bad request")
	}

	// Split request line
	splitRequest := strings.Split(splitHeaders[0], " ")
	if len(splitRequest) < 3 {
		conn.Write([]byte(GetStatus(400, "Bad Request")))
		return nil, fmt.Errorf("bad request")
	}

	// Extract headers
	headers := make(map[string]string)
	for _, header := range splitHeaders[1:] {
		parts := strings.SplitN(header, ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}

	// Get content encoding
	var contentEncoding string
	contentEncodingMethods := strings.Split(headers["Accept-Encoding"], ", ")
	for _, method := range contentEncodingMethods {
		fmt.Println(method)
		if method == "gzip" {
			contentEncoding = method
		}
	}

	req.Method = splitRequest[0]
	req.Path = splitRequest[1]
	req.Body = bodyPart
	req.UserAgent = headers["User-Agent"]
	req.ContentEncoding = contentEncoding

	return &req, nil
}
