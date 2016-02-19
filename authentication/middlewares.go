package authentication

type AuthMiddleware struct {
	AB backends.AuthBackend
}

func GetNewAuthenticationMiddleware(authBackend backends.AuthBackend) *AuthMiddleware {
	return &AuthMiddleware{AB: authBackend}
}
