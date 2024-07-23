package utils

import "fmt"

func GetStatus(code int, message string) string {
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", code, message)
}
