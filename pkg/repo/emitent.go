package repo

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
)

type EmitentRepository interface {
	Get() ([]models.User, error)
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(u *models.User) (*models.User, error)
}

type emitentRepository struct {
	db *gorm.DB
}

func NewEmitent(db *gorm.DB) EmitentRepository {
	return &emitentRepository{db}
}

func (r *emitentRepository) Get() ([]models.User, error) {
	var users []models.User

	if tx := r.db.Where(&models.User{Role: models.EmitentRole}).Find(&users); tx.Error != nil {
		return nil, tx.Error
	}

	return users, nil
}

func (r *emitentRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	res := r.db.Where(&models.User{ID: uint(id), Role: models.EmitentRole}).First(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (r *emitentRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	res := r.db.Where(&models.User{Email: email, Role: models.EmitentRole}).First(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (r *emitentRepository) Create(u *models.User) (*models.User, error) {
	if res := r.db.Create(&u); res.Error != nil {
		return nil, res.Error
	}

	return u, nil
}
