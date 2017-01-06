package main

import ()

type Simulate struct {
	hoverfly *Hoverfly
}

func (this Simulate) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	response, err := this.hoverfly.GetResponse(request, details)
	if err != nil {
		return hoverflyError(request, err, err.Error(), err.StatusCode), err
	}

	return response, nil
}