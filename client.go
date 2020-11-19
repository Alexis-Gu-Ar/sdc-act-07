package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Process struct {
	Id       int
	Progress int
}

func (p Process) print() {
	fmt.Printf("%d : %d\n", p.Id, p.Progress)
}

func (p *Process) increase() {
	p.Progress++
	time.Sleep(time.Millisecond * 500)
}

func requestProcess() *Process {
	var process *Process
	connection, _ := net.Dial("tcp", ":3000")
	gob.NewEncoder(connection).Encode("get process")
	gob.NewDecoder(connection).Decode(&process)
	connection.Close()
	return process
}

func returnProcessToServer(process *Process) {
	connection, _ := net.Dial("tcp", ":3000")

	var trasmiter *gob.Encoder = gob.NewEncoder(connection)
	trasmiter.Encode("return process")

	var receiver *gob.Decoder = gob.NewDecoder(connection)
	var response string
	receiver.Decode(&response)

	if response == "return it" {
		trasmiter.Encode(process)
	}
}

func startSesion() {
	var sigtermChan chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(sigtermChan, os.Interrupt, syscall.SIGTERM)

	var process *Process = requestProcess()
	for {
		select {
		case <-sigtermChan:
			returnProcessToServer(process)
			return
		default:
			process.print()
			process.increase()
		}
	}

}

func main() {
	startSesion()
}
