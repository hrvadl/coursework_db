package services

import (
	"context"
	"fmt"
	"time"

	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
)

type UserRepo interface {
	repo.Emitent
	repo.Stock
}

type Auth interface {
	SignUp(u *models.User) (*models.User, error)
	SignIn(u *models.User) (*models.User, error)
}

type auth struct {
	stock   repo.Stock
	emitent repo.Emitent
	crypto  Cryptor
	session repo.Session
}

func NewAuth(
	s repo.Stock,
	e repo.Emitent,
	se repo.Session,
	c Cryptor,
) Auth {
	return &auth{
		stock:   s,
		emitent: e,
		crypto:  c,
		session: se,
	}
}

func (a *auth) SignUp(u *models.User) (*models.User, error) {
	var userStrategy UserRepo

	switch u.Role {
	case models.StockRole:
		userStrategy = a.stock
	case models.EmitentRole:
		userStrategy = a.emitent
	default:
		return nil, fmt.Errorf("invalid role for auth: %v", u.Role)
	}

	if exists, _ := userStrategy.GetByEmail(u.Email); exists != nil {
		return nil, fmt.Errorf("user with email %v already exists", u.Email)
	}

	if u.Role != models.StockRole {
		return nil, fmt.Errorf("invalid role for stock: %v", u.Role)
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
	return userStrategy.Create(u)

}

func (a *auth) SignIn(u *models.User) (*models.User, error) {
	user, err := a.getUser(u.Email)

	if err != nil {
		return nil, err
	}

	if err := a.crypto.Compare(u.Password, user.Password); err != nil {
		return nil, fmt.Errorf("password does not match")
	}

	if _, err := a.session.Create(u); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *auth) getUser(email string) (*models.User, error) {
	emitentCh := make(chan *models.User, 1)
	stockCh := make(chan *models.User, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	go func() {
		if e, _ := a.emitent.GetByEmail(email); e != nil {
			emitentCh <- e
			return
		}
		emitentCh <- nil
	}()

	go func() {
		if s, _ := a.stock.GetByEmail(email); s != nil {
			stockCh <- s
			return
		}
		stockCh <- nil
	}()

	for {
		select {
		case u := <-stockCh:
			if u != nil {
				return u, nil
			}

		case u := <-emitentCh:
			if u != nil {
				return u, nil
			}

		case <-ctx.Done():
			return nil, fmt.Errorf("request timed out")
		}
	}
}
