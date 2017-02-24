package models

import "reflect"

type Simulation struct {
	Templates []RequestTemplateResponsePair
}

func NewSimulation() *Simulation {
	var templates []RequestTemplateResponsePair

	return &Simulation{
		Templates: templates,
	}
}

func (this *Simulation) AddRequestTemplateResponsePair(pair *RequestTemplateResponsePair) {
	var duplicate bool
	for _, savedPair := range this.Templates {
		duplicate = reflect.DeepEqual(pair.RequestTemplate, savedPair.RequestTemplate)
		if duplicate {
			break
		}
	}
	if !duplicate {
		this.Templates = append(this.Templates, *pair)
	}
}
