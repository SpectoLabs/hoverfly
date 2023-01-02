package models

import v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"

type Literal struct {
	Name  string
	Value interface{}
}

type Literals []Literal

func ImportLiterals(literals []v2.GlobalLiteralViewV5) *Literals {

	var allLiterals Literals
	for _, literal := range literals {
		allLiterals = append(allLiterals, Literal{
			Name:  literal.Name,
			Value: literal.Value,
		})
	}
	return &allLiterals
}

func (literals *Literals) ConvertToGlobalLiteralsPayloadView() []v2.GlobalLiteralViewV5 {

	var allLiterals []v2.GlobalLiteralViewV5
	if literals != nil {
		for _, literal := range *literals {

			allLiterals = append(allLiterals, v2.GlobalLiteralViewV5{
				Name:  literal.Name,
				Value: literal.Value,
			})

		}
	}
	return allLiterals
}
