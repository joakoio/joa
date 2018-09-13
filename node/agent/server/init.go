package server

import (
	zmq "github.com/alecthomas/gozmq"
	"github.com/joakoio/joa/common/conf"
	"sync"
)

type server struct {
	Socket *zmq.Socket
	Error  error
}

var instance *server
var once sync.Once

// server is declared as a singleton
//
// http://marcio.io/2015/07/singleton-pattern-in-go/
func GeInstance() *server {
	once.Do(func() {
		socket, err := initSocket()
		instance = &server{
			Socket: *socket,
			Error:  err,
		}
	})
	return instance
}

func initSocket() (*zmq.Socket, error) {

	context, _ := zmq.NewContext()
	//defer context.Close()

	socket, err := context.NewSocket(zmq.REP)
	socket.Bind("tcp://*:" + conf.GetAgentServerPort())

	return socket, err
}
