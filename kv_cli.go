package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/markbest/go-kv/app"
	"github.com/markbest/go-kv/utils/tcp"
)

var (
	inputAction string
	inputStr1   string
	inputStr2   string
	serverHost  = flag.String("h", "127.0.0.1", "tcp server host")
	serverPort  = flag.String("p", "9000", "tcp server port")
)

var help = func() {
	fmt.Println("USAGE: kv-cli -h [arguments] ...")
	fmt.Println("USAGE: kv-cli command [arguments] ...")
	fmt.Println("\nThe commands are:\n\taction\tkv [set|get|delete|list|persistent|exit|help]")
}

func main() {
	flag.Parse()

	// start client
	client := tcp.NewTCPClient(*serverHost, *serverPort, 10)
	for {
		fmt.Printf(">>> ")
		fmt.Scanln(&inputAction, &inputStr1, &inputStr2)
		switch inputAction {
		case "exit":
			rs, err := client.ReadWrite(func(conn *net.TCPConn) (string, error) {
				return app.HandleSendServerMsg(conn, "exit", inputStr1, inputStr2)
			})
			if err != nil {
				fmt.Printf("occur error %s", err.Error())
			}
			fmt.Println(rs)
			client.Close()
			goto Exit
		case "set":
			if inputStr1 == "" || inputStr2 == "" {
				help()
				continue
			}
			rs, err := client.ReadWrite(func(conn *net.TCPConn) (string, error) {
				return app.HandleSendServerMsg(conn, "set", inputStr1, inputStr2)
			})
			if err != nil {
				fmt.Printf("occur error %s", err.Error())
			}
			fmt.Println(rs)
		case "get":
			if inputStr1 == "" {
				help()
				continue
			}
			rs, err := client.ReadWrite(func(conn *net.TCPConn) (string, error) {
				return app.HandleSendServerMsg(conn, "get", inputStr1, inputStr2)
			})
			if err != nil {
				fmt.Printf("occur error %s", err.Error())
			}
			fmt.Println(rs)
		case "delete":
			if inputStr1 == "" {
				help()
				continue
			}
			rs, err := client.ReadWrite(func(conn *net.TCPConn) (string, error) {
				return app.HandleSendServerMsg(conn, "delete", inputStr1, inputStr2)
			})
			if err != nil {
				fmt.Printf("occur error %s", err.Error())
			}
			fmt.Println(rs)
		case "list":
			rs, err := client.ReadWrite(func(conn *net.TCPConn) (string, error) {
				return app.HandleSendServerMsg(conn, "list", inputStr1, inputStr2)
			})
			if err != nil {
				fmt.Printf("occur error %s", err.Error())
			}
			fmt.Println(rs)
		case "persistent":
			rs, err := client.ReadWrite(func(conn *net.TCPConn) (string, error) {
				return app.HandleSendServerMsg(conn, "persistent", inputStr1, inputStr2)
			})
			if err != nil {
				fmt.Printf("occur error %s", err.Error())
			}
			fmt.Println(rs)
		case "help":
			help()
			continue
		default:
			help()
			continue
		}
	}
Exit:
}
