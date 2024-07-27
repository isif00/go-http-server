package handlers

import "net"

func ConnectionHandler(conn net.Conn) {
	defer conn.Close()
	for {
		if err := RequestHandler(conn); err != nil {
			break
		}
	}
}
