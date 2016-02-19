package authentication

type JWTAuthenticationBackend struct {
	privateKey  *rsa.PrivateKey
	PublicKey   *rsa.PublicKey
	AuthBackend backends.AuthBackend
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

var authBackendInstance *JWTAuthenticationBackend = nil

// Token - container for jwt.Token for encoding
type Token struct {
	Token *jwt.Token
}

func (t *Token) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(t)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeToken(data []byte) (*Token, error) {
	var t *Token
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

