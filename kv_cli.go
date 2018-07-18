package main

import (
	"flag"
	"fmt"

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
			rs := app.HandleScanInput(client, "exit", inputStr1, inputStr2)
			fmt.Println(rs)
			client.Close()
			goto Exit
		case "set":
			if inputStr1 == "" || inputStr2 == "" {
				help()
				continue
			}
			rs := app.HandleScanInput(client, "set", inputStr1, inputStr2)
			fmt.Println(rs)
			app.ClearScan(&inputStr1, &inputStr2)
		case "get":
			if inputStr1 == "" {
				help()
				continue
			}
			rs := app.HandleScanInput(client, "get", inputStr1, inputStr2)
			fmt.Println(rs)
			app.ClearScan(&inputStr1, &inputStr2)
		case "delete":
			if inputStr1 == "" {
				help()
				continue
			}
			rs := app.HandleScanInput(client, "delete", inputStr1, inputStr2)
			fmt.Println(rs)
			app.ClearScan(&inputStr1, &inputStr2)
		case "list":
			rs := app.HandleScanInput(client, "list", inputStr1, inputStr2)
			fmt.Println(rs)
			app.ClearScan(&inputStr1, &inputStr2)
		case "persistent":
			rs := app.HandleScanInput(client, "persistent", inputStr1, inputStr2)
			fmt.Println(rs)
			app.ClearScan(&inputStr1, &inputStr2)
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
