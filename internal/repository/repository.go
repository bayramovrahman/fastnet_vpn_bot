package repository

import "github.com/bayramovrahman/fastnet_vpn_bot/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	GetUserById(id int) (models.User, error)
	UpdateUser(user models.User) error
	Authenticate(email, testPassword string) (int, string, error)

	// User Login Security methods
	GetUserLoginSecurity(userID int) (models.UserLoginSecurity, error)
	UpdateUserLoginSecurity(security models.UserLoginSecurity) error
	CreateUserLoginSecurity(security models.UserLoginSecurity) error
}
