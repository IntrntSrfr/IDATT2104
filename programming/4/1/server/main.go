package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {

	ln, err := net.ListenPacket("udp", ":54321")
	if err != nil {
		return
	}
	defer ln.Close()

	for {
		tmp := make([]byte, 128)
		n, addr, err := ln.ReadFrom(tmp)
		if err != nil {
			continue
		}
		handle(ln, tmp[:n], n, addr)
	}
}

var FailErr = []byte("illegal expression - make sure there are two numbers and plus or minus")

func handle(c net.PacketConn, data []byte, n int, addr net.Addr) {
	exprStr := string(data)

	if exprStr == "" {
		return
	}

	//fmt.Println(exprStr, n, []byte(exprStr))
	fmt.Println(addr, n)

	if tmp := strings.Split(exprStr, "+"); len(tmp) == 2 {
		nums, err := splitNums(tmp)
		if err != nil {
			c.WriteTo(FailErr, addr)
			return
		}
		c.WriteTo([]byte(fmt.Sprint(nums[0]+nums[1])), addr)
	} else if tmp := strings.Split(exprStr, "-"); len(tmp) == 2 {
		nums, err := splitNums(tmp)
		if err != nil {
			c.WriteTo(FailErr, addr)
			return
		}
		c.WriteTo([]byte(fmt.Sprint(nums[0]-nums[1])), addr)
	} else {
		c.WriteTo(FailErr, addr)
	}
}

func splitNums(nums []string) ([]int, error) {
	num1, err := strconv.Atoi(nums[0])
	if err != nil {
		return nil, err
	}

	num2, err := strconv.Atoi(nums[1])
	if err != nil {
		return nil, err
	}
	return []int{num1, num2}, nil
}
