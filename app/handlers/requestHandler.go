package handlers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/config"
	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

func RequestHandler(conn net.Conn) error {
	// Read the request
	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		fmt.Println("Error reading request:", err.Error())
		return err
	}

	requestLine := string(req[:n])
	fmt.Printf("Request: %s\n", requestLine)

	parsedRequest, nil := utils.ParseRequest(requestLine, conn)
	if nil != nil {
		conn.Write([]byte(utils.GetStatus(400, "Bad Request")))
		return err
	}

	var response string

	var method = parsedRequest.Method
	var body = parsedRequest.Body
	var contentEncoding = parsedRequest.ContentEncoding

	switch path := parsedRequest.Path; {

	case path == "/":
		response = utils.GetStatus(200, "OK\r\n")

	case strings.HasPrefix(path, "/echo/"):
		parts := strings.SplitN(path, "/", 3)
		if len(parts) < 3 {
			conn.Write([]byte(utils.GetStatus(400, "Bad Request\r\n")))
			return err
		}
		message := parts[2]
		if contentEncoding == "gzip" {
			var buf bytes.Buffer
			w := gzip.NewWriter(&buf)
			_, err := w.Write([]byte(message))
			if err != nil {
				conn.Write([]byte(utils.GetStatus(500, "Internal Server Error\r\n")))
				return err
			}
			w.Close()
			response = fmt.Sprintf("%sContent-Type: text/plain\r\nContent-Encoding: gzip\r\nContent-Length: %d\r\n\r\n%s",
				utils.GetStatus(200, "OK"), buf.Len(), buf.Bytes())
		} else {
			response = fmt.Sprintf("%sContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
				utils.GetStatus(200, "OK"), len(message), message)
		}

	case strings.HasPrefix(path, "/user-agent"):
		response = fmt.Sprintf("%sContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", utils.GetStatus(200, "OK"), len(parsedRequest.UserAgent), parsedRequest.UserAgent)

	case strings.HasPrefix(path, "/files/"):
		fileName := strings.TrimPrefix(path, "/files/")
		dir := config.Directory
		if method == "GET" {
			data, err := os.ReadFile(dir + "/" + fileName)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				conn.Write([]byte(utils.GetStatus(404, "Not Found\r\n")))
				return err
			}

			response = fmt.Sprintf("%sContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", utils.GetStatus(200, "OK"), len(data), string(data))

		} else if method == "POST" {
			file := []byte(body)
			if err := os.WriteFile(dir+"/"+fileName, file, 0644); err == nil {
				response = utils.GetStatus(201, "Created\r\n")
			}
		}

	default:
		response = utils.GetStatus(404, "Not Found\r\n")
	}

	conn.Write([]byte(response))
	conn.Close()
	return nil
}
