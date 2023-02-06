package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		return
	}
	conn, err := listener.Accept()
	if err != nil {
		return
	}
	defer conn.Close()

	req := make([]byte, 1024)
	conn.Read(req)
	w := strings.Split(strings.TrimSpace(string(req)), "\r\n")
	list := ""
	for _, s := range w[:len(w)-2] {
		list += fmt.Sprintf("<li>%v</li>", s)
	}

	sb := strings.Builder{}

	sb.Write([]byte("HTTP/1.0 200 OK\r\n"))
	sb.Write([]byte("Content-Type: text/html; charset=utf-8\r\n\r\n"))
	sb.Write([]byte("<html><body>"))
	sb.Write([]byte("<h1> Dette er en fantastisk webserver </h1>"))
	sb.Write([]byte("<p> Her er headerne fra klienten: </p>"))
	sb.Write([]byte("<ul>"))
	sb.Write([]byte(list))
	sb.Write([]byte("</ul>"))
	sb.Write([]byte("</body></html>"))

	conn.Write([]byte(sb.String()))
}
