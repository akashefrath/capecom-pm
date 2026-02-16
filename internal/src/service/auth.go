package service

import "github.com/akashefrath/capecom-pm/internal/src/repository"

type Auth struct {
	Auth *repository.Auth
	User *repository.User
}

func NewAuth(auth *repository.Auth, user *repository.User) *Auth {
	return &Auth{
		Auth: auth,
		User: user,
	}
}

func (a *Auth) Login(email, password string) error {
	return a.User.GetUserByEmail(email)
}
