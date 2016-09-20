package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/core/authentication"
	"github.com/go-zoo/bone"
)

type AdminHandler interface {
	RegisterRoutes(*bone.Mux, *authentication.AuthMiddleware)
}
