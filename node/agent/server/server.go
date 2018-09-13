package server

import (
	"log"
	"github.com/joakoio/joa/node/agent/subscriber"
)

func (s *server) StartListening(){
	defer s.Socket.Close()

	sub := subscriber.GeInstance()
	for {
		message, err := s.Socket.Recv(0)
		s.Error = err
		log.Println(string(message))

		result := "Success"
		if err != nil {
			log.Println(err)
			result = err.Error()
		}

		sub.AsyncBrodcast(string(message))

		s.Socket.Send([]byte(result),0)
	}
}
