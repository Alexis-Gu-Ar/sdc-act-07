package main

import "fmt"

func main() {
	var myServer = Server(3000, 3)
	go myServer.Start()
	fmt.Scanln()
}
