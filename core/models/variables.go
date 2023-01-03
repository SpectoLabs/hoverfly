package models

import (
	"fmt"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

type Variable struct {
	Name      string
	Function  string
	Arguments []interface{}
}

type Variables []Variable

func ValidateVariablePayload(variables []v2.GlobalVariableViewV5, helperMethodMap map[string]interface{}) (err error) {

	for _, variable := range variables {
		if _, found := helperMethodMap[variable.Function]; !found {

			return fmt.Errorf("function %s not supported for custom variable %s", variable.Function, variable.Name)
		}
	}
	return nil
}

func ImportVariables(variables []v2.GlobalVariableViewV5) *Variables {

	var allVariables Variables
	for _, variable := range variables {
		allVariables = append(allVariables, Variable{
			Name:      variable.Name,
			Function:  variable.Function,
			Arguments: variable.Arguments,
		})
	}
	return &allVariables
}

func (variables *Variables) ConvertToGlobalVariablesPayloadView() []v2.GlobalVariableViewV5 {

	var allVariablesView []v2.GlobalVariableViewV5
	if variables != nil {
		for _, variable := range *variables {
			variableView := v2.GlobalVariableViewV5{
				Name:      variable.Name,
				Function:  variable.Function,
				Arguments: variable.Arguments,
			}
			allVariablesView = append(allVariablesView, variableView)

		}
	}

	return allVariablesView
}
