package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("please provide host:port")
	}

	fmt.Println(os.Args)
	addr := os.Args[1]

	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		//text = text[:len(text)-2]

		if strings.TrimSpace(text) == "exit" {
			break
		}
		if text == "" {
			continue
		}

		fmt.Fprintf(conn, text)

		conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		msgByte := make([]byte, 128)
		_, err := conn.Read(msgByte)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Print("-> " + string(msgByte))
		fmt.Println("")
	}
}
