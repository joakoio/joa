package main

import (
	"os"
	"flag"
	"fmt"
	"time"
	zmq "github.com/alecthomas/gozmq"
)

const (
	PORT_FRONT   = 9001
	PORT_BACK    = 5559
	PORT_MONITOR = 9003
)

type Goque struct{
	ctx *zmq.Context
}

func NewGoque() *Goque {
	ctx, _ := zmq.NewContext()
	defer ctx.Close()

	return &Goque{ctx}
}

func (goque *Goque) Client() {
	sock, _:= goque.ctx.NewSocket(zmq.REQ)
	defer sock.Close()

	sock.Connect(fmt.Sprintf("tcp://localhost:%d", PORT_FRONT))
	fmt.Printf("start client...\n")
	for {
		sendmsg := fmt.Sprintf("PING#<%d>", os.Getpid())
		sock.Send([]byte(sendmsg), 0)
		msg, _ := sock.Recv(0)
		fmt.Println("Recv:", string(msg))
		time.Sleep(1 * time.Second)
	}
}

func (goque *Goque) Server() {
	sock, _:= goque.ctx.NewSocket(zmq.REP)
	defer sock.Close()

	fmt.Println(fmt.Sprintf("tcp://localhost:%d", PORT_BACK))
	sock.Connect(string(fmt.Sprintf("tcp://localhost:%d", PORT_BACK)))
	fmt.Printf("start server...\n")
	for {
		msg, _ := sock.Recv(0)
		fmt.Println("Recv:", string(msg))
		sendmsg := fmt.Sprintf("PONG#<%d>", os.Getpid())
		sock.Send([]byte(sendmsg), 0)
	}
}

func (goque *Goque) Monitor() {
	sock, _:= goque.ctx.NewSocket(zmq.SUB)
	defer sock.Close()

	sock.Connect(fmt.Sprintf("tcp://localhost:%d", PORT_MONITOR))
	sock.SetSubscribe("")
	for {
		msg, _ := sock.Recv(0)
		fmt.Println("MONITOR:", string(msg))
	}
}

func (goque *Goque) Queue() {
	front, _ := goque.ctx.NewSocket(zmq.ROUTER)
	defer front.Close()
	front.Bind(fmt.Sprintf("tcp://*:%d", PORT_FRONT))

	back, _ := goque.ctx.NewSocket(zmq.DEALER)
	defer back.Close()
	back.Bind(fmt.Sprintf("tcp://*:%d", PORT_BACK))

	zmq.Device(zmq.QUEUE, front, back)
}

func (goque *Goque) MonitoredQueue() {
	front, _ := goque.ctx.NewSocket(zmq.ROUTER)
	defer front.Close()
	front.Bind(fmt.Sprintf("tcp://*:%d", PORT_FRONT))

	back, _ := goque.ctx.NewSocket(zmq.DEALER)
	defer back.Close()
	back.Bind(fmt.Sprintf("tcp://*:%d", PORT_BACK))

	mon, _ := goque.ctx.NewSocket(zmq.PUB)
	defer mon.Close()
	mon.Bind(fmt.Sprintf("tcp://*:%d", PORT_MONITOR))

	polls := zmq.PollItems {
		//zmq.PollItem{Socket: front, a: zmq.POLLIN},
		//zmq.PollItem{Socket: back , zmq.Events: zmq.POLLIN},
	}

	total := make(map[string]int)
	for {
		_, _ = zmq.Poll(polls, -1)
		switch {
		case polls[0].REvents & zmq.POLLIN != 0:
			parts, _ := front.RecvMultipart(0)
			back.SendMultipart(parts, 0)
			mon.Send([]byte(fmt.Sprintf("IN: %d, OUT %d", total["in"], total["out"])), 0)
			total["in"] += 1
		case polls[1].REvents & zmq.POLLIN != 0:
			parts, _ := back.RecvMultipart(0)
			front.SendMultipart(parts, 0)
			total["out"] += 1
		}
	}
}

func (goque *Goque) Run(key string) {
	switch key {
	case "serv":
		goque.Server()
	case "client":
		goque.Client()
	case "queue":
		goque.Queue()
	case "monitorq":
		goque.MonitoredQueue()
	case "monitor":
		goque.Monitor()
	default:
		fmt.Printf("serv or client\n")
	}
}


func usage() {
	fmt.Println("Usage: go run queue.go [serv|client|queue|monitorq|monitor]\n")
	flag.PrintDefaults()
	os.Exit(2)
}


func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	goque := NewGoque()
	goque.Run(args[0])

	os.Exit(0)
}

//go run queue.go server
//go run queue.go monitorq
//go run queue.go monitor
//python queue.go client