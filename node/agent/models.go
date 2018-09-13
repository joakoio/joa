package agent

import "time"

type Agent struct{
	Name string
	Description string
	Namespace string
	Status Status
}


type Status struct {
	Active bool
	LastCommunication time.Time
}