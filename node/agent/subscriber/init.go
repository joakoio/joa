package subscriber

import (
	zmq "github.com/alecthomas/gozmq"
	"github.com/joakoio/joa/common/conf"
	"github.com/joakoio/joa/node/agent"
	"sync"
)

type subscriber struct {
	Agents []agent.Agent
	Socket *zmq.Socket
	Error  error
}

var instance *subscriber
var once sync.Once

// subscriber is declared as a singleton
//
// http://marcio.io/2015/07/singleton-pattern-in-go/
func GeInstance() *subscriber {
	once.Do(func() {
		socket, err := initSocket()
		instance = &subscriber{
			Socket: *socket,
			Error:  err,
		}
	})
	return instance
}

func initSocket() (*zmq.Socket, error) {
	context, _ := zmq.NewContext()
	//defer context.Close()

	// This is our public endpoint for subscribers
	socket, err := context.NewSocket(zmq.PUB)
	socket.Bind("tcp://*:" + conf.GetAgentSubscriberPort())

	return socket, err
}
