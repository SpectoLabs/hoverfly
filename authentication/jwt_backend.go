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

func InitJWTAuthenticationBackend(ab backends.AuthBackend) *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey:  getPrivateKey(),
			PublicKey:   getPublicKey(),
			AuthBackend: ab,
		}
	}

	return authBackendInstance
}

func (backend *JWTAuthenticationBackend) GenerateToken(userUUID string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(Get().JWTExpirationDelta)).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["sub"] = userUUID
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("got error while generating JWT token")
		return "", err
	}
	return tokenString, nil
}

func (backend *JWTAuthenticationBackend) Authenticate(user *User) bool {
	// TODO: get user info
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testing"), 10)

	testUser := User{
		UUID:     uuid.New(),
		Username: "hoverfly",
		Password: string(hashedPassword),
	}

	return user.Username == testUser.Username && bcrypt.CompareHashAndPassword([]byte(testUser.Password), []byte(user.Password)) == nil
}

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	return backend.AuthBackend.Delete([]byte(tokenString))
}

func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {

	redisToken, _ := backend.AuthBackend.GetValue([]byte(token))

	if redisToken == nil {
		return false
	}
	return true
}
