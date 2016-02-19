package controllers

type AuthController struct {
	AB backends.AuthBackend
}

func GetNewAuthenticationController(authBackend backends.AuthBackend) *AuthController {
	return &AuthController{AB: authBackend}
}
