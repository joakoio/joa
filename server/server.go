package server

import (
	"fmt"
	"os"

	zmq "github.com/alecthomas/gozmq"
	"time"
)

const (
	PORT = 5559
)

type Goque struct {
	ctx *zmq.Context
}

func (goque *Goque) Server() {
	//sock, _ := goque.ctx.NewSocket(zmq.REP)
	//defer sock.Close()
	//
	//fmt.Println(fmt.Sprintf("tcp://localhost:%d", PORT))
	//sock.Connect(string(fmt.Sprintf("tcp://localhost:%d", PORT)))
	//fmt.Printf("start server...\n")
	//for {
	//	msg, _ := sock.Recv(0)
	//	fmt.Println("Recv:", string(msg))
	//	sendmsg := fmt.Sprintf("PONG#<%d>", os.Getpid())
	//	sock.Send([]byte(sendmsg), 0)
	//}
	context, _ := zmq.NewContext()
	socket, _ := context.NewSocket(zmq.REP)
	defer context.Close()
	defer socket.Close()
	socket.Bind("tcp://*:5555")

	// Wait for messages
	for {
		fmt.Printf("Waiting for message")
		msg, _ := socket.Recv(0)
		println("Received ", string(msg))

		// do some fake "work"
		//time.Sleep(3* time.Millisecond)
		time.Sleep(time.Second)

		// send reply back to client
		reply := fmt.Sprintf("World")
		socket.Send([]byte(reply), 0)
	}
}

func (goque *Goque) Client() {
	sock, _ := goque.ctx.NewSocket(zmq.REQ)
	defer sock.Close()

	sock.Connect(fmt.Sprintf("tcp://localhost:%d", PORT))
	fmt.Printf("start client...\n")
	for {
		sendmsg := fmt.Sprintf("PING#<%d>", os.Getpid())
		sock.Send([]byte(sendmsg), 0)
		msg, _ := sock.Recv(0)
		fmt.Println("Recv:", string(msg))
		time.Sleep(1 * time.Second)
	}
}
