package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("please provide host:port")
	}

	fmt.Println(os.Args)
	addr := os.Args[1]

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	msgByte := make([]byte, 1024)
	conn.Read(msgByte)
	fmt.Print("->: " + strings.TrimSpace(string(msgByte)))

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')

		if strings.TrimSpace(string(text)) == "exit" {
			fmt.Fprintf(conn, text+"\n")
			break
		}

		fmt.Fprintf(conn, text+"\n")

		msgByte = make([]byte, 1024)
		conn.Read(msgByte)
		fmt.Print("->: " + string(msgByte))
		fmt.Println("")
	}
}
