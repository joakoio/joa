package subscriber

func (s *subscriber) Brodcast(message string){
	//message := "This is a message"
	m := newMessage("A", string(message))
	s.Socket.SendMultipart(m, 0)
	//backend.Send(message, 0)
}

func (s *subscriber) AsyncBrodcast(message string){
	go s.Brodcast(message)
}


func newMessage(filter, message string) [][]byte {
	return [][]byte{[]byte("A"), []byte(message)}
}