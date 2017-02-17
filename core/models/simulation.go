package models

type Simulation struct {
	Templates []RequestResponsePair
}

func NewSimulation() *Simulation {
	var templates []RequestResponsePair

	return &Simulation{
		Templates: templates,
	}
}
