package main

import "fmt"

func main() {
	totalProcess := 5
	var myServer = Server(3000, totalProcess)
	go myServer.Start()
	fmt.Scanln()
}
