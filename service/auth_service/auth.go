package auth_service

import (
	"oamp/models"
)

type Auth struct {
	Username string
	Password string
}

func (a *Auth) ExistsUserName() (bool, error) {
	return models.ExistsUserName(a.Username)
}
func (a *Auth) ExistsPassword() (bool, error) {
	return models.CheckPassword(a.Username, a.Password)
}
