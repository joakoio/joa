package main

import (
	"fmt"
	"time"
	zmq "github.com/alecthomas/gozmq"
)

func main () {
	context, _ := zmq.NewContext()
	socket, _ := context.NewSocket(zmq.REQ)
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