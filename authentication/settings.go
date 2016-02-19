package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var environments = map[string]string{
	"production":    "../../authentication/settings/prod.json",
	"preproduction": "../../authentication/settings/pre.json",
	"tests":         "../../authentication/settings/tests.json",
}

type Settings struct {
	PrivateKeyPath     string
	PublicKeyPath      string
	JWTExpirationDelta int
}

