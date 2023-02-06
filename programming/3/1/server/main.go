package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {

	ln, err := net.Listen("tcp", ":54321")
	if err != nil {
		return
	}
	defer ln.Close()

	fmt.Println(ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handle(conn)
	}
}

func handle(c net.Conn) {
	defer c.Close()

	fmt.Println("serving new connection:", c.RemoteAddr())
	fmt.Fprint(c, "welcome to the calculator! - type exit to leave\nplease input your expression:\n")
	for {

		exprStr, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println("fail army 1")
			break
		}

		exprStr = strings.TrimSpace(string(exprStr))

		if exprStr == "exit" {
			fmt.Println("closing connection", c.RemoteAddr())
			break
		}

		if exprStr == "" {
			fmt.Fprint(c, " ")
			continue
		}

		fmt.Println(exprStr)

		if tmp := strings.Split(exprStr, "+"); len(tmp) == 2 {
			fmt.Println("plus expression")
			nums, err := splitNums(tmp)
			if err != nil {
				fmt.Fprint(c, "illegal expression - make sure there are two numbers and plus or minus")
				continue
			}
			fmt.Fprint(c, nums[0]+nums[1])
		} else if tmp := strings.Split(exprStr, "-"); len(tmp) == 2 {
			fmt.Println("minus expression")
			nums, err := splitNums(tmp)
			if err != nil {
				fmt.Fprint(c, "illegal expression - make sure there are two numbers and plus or minus")
				continue
			}
			fmt.Fprint(c, nums[0]-nums[1])
		} else {
			fmt.Fprint(c, "illegal expression - make sure there are two numbers and plus or minus")
		}
	}
}

func splitNums(nums []string) ([]int, error) {
	num1, err := strconv.Atoi(nums[0])
	if err != nil {
		fmt.Println("fail army 3")
		return nil, err
	}

	num2, err := strconv.Atoi(nums[1])
	if err != nil {
		fmt.Println("fail army 4")
		return nil, err
	}
	return []int{num1, num2}, nil
}
