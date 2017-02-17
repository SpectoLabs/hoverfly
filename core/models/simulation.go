package models

type Simulation struct {
	Templates []RequestTemplateResponsePair
}

func NewSimulation() *Simulation {
	var templates []RequestTemplateResponsePair

	return &Simulation{
		Templates: templates,
	}
}
