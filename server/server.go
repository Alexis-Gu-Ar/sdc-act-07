package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"time"
)

type ProcessFactory struct {
	currentId int
}

func (factory *ProcessFactory) getProcess() Process {
	defer func() {
		factory.currentId++
	}()
	return Process{
		Id: factory.currentId,
	}
}

type Process struct {
	Id              int
	Progress        int
	runningOnClient bool
}

func (p Process) print() {
	fmt.Printf("ID: %d, progress: %d\n", p.Id, p.Progress)
}

type server struct {
	port               int
	processesDelegated int
	processes          []Process
	processFactory     ProcessFactory
}

func Server(port int, totalProcesses int) server {
	s := server{
		port: port,
	}

	for h := 0; h < totalProcesses; h++ {
		s.processes = append(s.processes, s.processFactory.getProcess())
	}

	return s
}

func (s *server) Start() {
	go s.Listen()
	for {
		for i := 0; i < len(s.processes); i++ {
			if !s.processes[i].runningOnClient {
				s.processes[i].print()
				s.processes[i].Progress++
			}
		}
		fmt.Println("------------------------------")
		time.Sleep(time.Millisecond * 500)
	}
}

func (s *server) HandleConnection(connection net.Conn) {
	var tx *gob.Encoder = gob.NewEncoder(connection)
	var rx *gob.Decoder = gob.NewDecoder(connection)

	var request string
	rx.Decode(&request)
	if request == "get process" {
		for i := 0; i < len(s.processes); i++ {
			if !s.processes[i].runningOnClient {
				s.processes[i].runningOnClient = true
				s.processesDelegated++
				gob.NewEncoder(connection).Encode(s.processes[i])
				break
			}
		}
	} else if request == "return process" {
		tx.Encode("return it")
		var delegatedProcess *Process
		rx.Decode(&delegatedProcess)

		for i := 0; i < len(s.processes); i++ {
			if delegatedProcess.Id == s.processes[i].Id {
				s.processes[i].Progress = delegatedProcess.Progress
				s.processes[i].runningOnClient = false
				break
			}
		}
	}
	connection.Close()
}

func (s *server) Listen() {
	listener, _ := net.Listen("tcp", ":"+strconv.Itoa(s.port))
	fmt.Printf("Listening at port %d\n", s.port)

	for {
		connection, _ := listener.Accept()
		go s.HandleConnection(connection)
	}
}
