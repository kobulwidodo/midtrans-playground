package http

import "go-midtrans/domain"

type authResponse struct {
	Name  string
	Email string
	Token string
}

func AuthResponse(user domain.User, token string) authResponse {
	return authResponse{
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}
}
