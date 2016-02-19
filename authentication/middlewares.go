package authentication

type AuthMiddleware struct {
	AB backends.AuthBackend
}
