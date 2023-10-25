package services

import (
	"fmt"

	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type Auth interface {
	SignUp(u *models.User) (*models.User, error)
	SignIn(u *models.User) (*models.User, *models.Session, error)
}

type auth struct {
	user    repo.User
	crypto  Cryptor
	session repo.Session
}

func NewAuth(
	u repo.User,
	se repo.Session,
	c Cryptor,
) Auth {
	return &auth{
		user:    u,
		crypto:  c,
		session: se,
	}
}

func (a *auth) SignUp(u *models.User) (*models.User, error) {
	if exists, _ := a.user.GetByEmail(u.Email); exists != nil {
		return nil, fmt.Errorf("user with email %v already exists", u.Email)
	}

	if u.Balance < 0 {
		return nil, fmt.Errorf("balance cannot be negative")
	}

	if u.FirstName == "" {
		return nil, fmt.Errorf("first name cannot be empty")
	}

	if u.LastName == "" {
		return nil, fmt.Errorf("last name cannot be empty")
	}

	if len(u.Password) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters long")
	}

	pass, err := a.crypto.Hash(u.Password)

	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %v", err)
	}

	u.Password = pass
	return a.user.Create(u)

}

func (a *auth) SignIn(u *models.User) (*models.User, *models.Session, error) {
	user, err := a.user.GetByEmail(u.Email)

	if err != nil {
		return nil, nil, err
	}

	if err := a.crypto.Compare(u.Password, user.Password); err != nil {
		return nil, nil, fmt.Errorf("password does not match")
	}

	sess, err := a.session.Create(user)

	if err != nil {
		return nil, nil, err
	}

	return user, sess, nil
}
