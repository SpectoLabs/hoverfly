package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

type ClosestMiss struct {
	RequestDetails RequestDetails
	Response       v2.ResponseDetailsViewV5
	RequestMatcher v2.RequestMatcherViewV5
	MissedFields   []string
	State          map[string]string
}

func (this *ClosestMiss) GetMessage() string {

	requestBytes, _ := json.MarshalIndent(this.RequestDetails, "", "    ")
	matcherBytes, _ := json.MarshalIndent(this.RequestMatcher, "", "    ")
	responseBytes, _ := json.MarshalIndent(this.Response, "", "    ")
	currentState, _ := json.MarshalIndent(this.State, "", "    ")

	return "\n\nThe following request was made, but was not matched by Hoverfly:\n\n" +
		string(requestBytes) +
		"\n\nWhilst Hoverfly has the following state:\n\n" +
		string(currentState) +
		"\n\nThe matcher which came closest was:\n\n" +
		string(matcherBytes) +
		"\n\nBut it did not match on the following fields:\n\n" +
		fmt.Sprint("["+strings.Join(this.MissedFields, ", ")+"]") +
		"\n\nWhich if hit would have given the following response:\n\n" +
		string(responseBytes)
}

func (this *ClosestMiss) BuildView() *v2.ClosestMissView {
	return &v2.ClosestMissView{
		Response:       this.Response,
		RequestMatcher: this.RequestMatcher,
		MissedFields:   this.MissedFields,
	}
}
