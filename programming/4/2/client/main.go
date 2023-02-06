package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("please provide host:port")
	}

	addr := os.Args[1]

	rootCAs := x509.NewCertPool()
	cert, err := ioutil.ReadFile("server.crt")
	if err != nil {
		log.Fatal(err)
	}

	if ok := rootCAs.AppendCertsFromPEM(cert); !ok {
		log.Fatal("No certs appended, using system certs only")
	}

	conf := &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: false,
	}

	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	msgByte := make([]byte, 1024)
	conn.Read(msgByte)
	fmt.Print("->: " + string(msgByte))

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
