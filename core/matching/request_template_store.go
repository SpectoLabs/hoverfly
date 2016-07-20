package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	"errors"
)

type RequestTemplateStore struct {

}

func(this *RequestTemplateStore) GetPayload() (*models.Payload, error) {
	return nil, errors.New("not implemented")
}