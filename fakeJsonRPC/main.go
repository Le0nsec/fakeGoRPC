package main

/*
* @Description: A fake Golang jsonRPC service.
* @Author: Leonsec
* @Date: 2022-03-20
* @LastEditTime: 2022-03-20
 */

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// var returnMsg = `{"id":<ID>,"result":{"Pro":666,"Quo":0,"Rem":0},"error":null}`

var (
	port      int
	returnMsg string
)

func init() {
	flag.IntVar(&port, "p", 8086, "listen `port`")
	flag.StringVar(&returnMsg, "re", "", "`json` message that will be return")
}

func main() {
	flag.Parse()
	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("listen error: ", err)
	}
	defer listener.Close()
	log.Printf("listen at %s", address)

	if returnMsg == "" {
		fmt.Println("[*] You have not specified the returned message, use `-h` to view the usage")
	} else {
		fmt.Printf("[*] The return message is %s\n", returnMsg)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept connection error: %s", err.Error())
		}
		go Handler(conn)
	}
}

func Handler(conn net.Conn) {
	defer conn.Close()
	log.Printf("new connection from %s", conn.RemoteAddr().String())

	buf := make([]byte, 512)
	len, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("[!] Receive data error: %s\n", err.Error())
		return
	}
	if len == 0 {
		return
	}
	buf = buf[:len]
	fmt.Printf("[+] Receive data: %s", buf)

	handlerConn(&conn, buf)
}

func handlerConn(conn *net.Conn, buf []byte) {
	var m map[string]interface{}
	if json.Valid(buf) {
		err := json.Unmarshal(buf, &m)
		if err != nil {
			fmt.Printf("[!] Unmarshal json error: %s", err.Error())
		}
	} else {
		fmt.Println("[!] Json format invalid!")
		return
	}
	if returnMsg != "" {
		msg := strings.Replace(returnMsg, "<ID>", strconv.Itoa(int(m["id"].(float64))), -1)
		writeConn(conn, []byte(msg))
	}
	buf = make([]byte, 512)
	len, err := (*conn).Read(buf)
	if err != nil {
		if err.Error() == "EOF" {
			return
		}
		fmt.Printf("[!] Receive data error: %s\n", err.Error())
		return
	}
	if len == 0 {
		return
	}
	buf = buf[:len]
	fmt.Printf("[+] Receive data: %s", buf)

	handlerConn(conn, buf)
}

func writeConn(conn *net.Conn, msg []byte) {
	_, err := (*conn).Write(msg)
	if err != nil {
		fmt.Printf("[!] Write data error: %s\n", err.Error())
		if err.Error() == "EOF" {
			(*conn).Close()
			return
		}
	}
	fmt.Printf("[-] Return data: %s\n", msg)
}
