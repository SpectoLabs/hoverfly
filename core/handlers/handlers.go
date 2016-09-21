package handlers

import (
	"github.com/go-zoo/bone"
)

type AdminHandler interface {
	RegisterRoutes(*bone.Mux, *AuthHandler)
}
