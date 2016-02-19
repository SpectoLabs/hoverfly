package authentication

import (
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/authentication/backends"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

type TokenAuthentication struct {
	Token string `json:"token" form:"token"`
}

type User struct {
	UUID     string `json:"uuid" form:"-"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
