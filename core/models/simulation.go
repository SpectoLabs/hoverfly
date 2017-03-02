package models

import (
	"reflect"
)

type Simulation struct {
	Templates      []RequestTemplateResponsePair
	ResponseDelays ResponseDelays
}

func NewSimulation() *Simulation {

	return &Simulation{
		Templates:      []RequestTemplateResponsePair{},
		ResponseDelays: &ResponseDelayList{},
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
