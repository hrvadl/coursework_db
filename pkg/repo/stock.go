package repo

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
)

type StockRepository interface {
	Get() ([]models.User, error)
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(u *models.User) (*models.User, error)
}

type stockRepository struct {
	db *gorm.DB
}

func NewStock(db *gorm.DB) StockRepository {
	return &stockRepository{db}
}

func (r *stockRepository) Get() ([]models.User, error) {
	var users []models.User

	if tx := r.db.Where(&models.User{Role: models.StockRole}).Find(&users); tx.Error != nil {
		return nil, tx.Error
	}

	return users, nil
}

func (r *stockRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	res := r.db.Where(&models.User{ID: uint(id), Role: models.StockRole}).First(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (r *stockRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	res := r.db.Where(&models.User{Email: email, Role: models.StockRole}).First(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (r *stockRepository) Create(u *models.User) (*models.User, error) {
	if res := r.db.Create(&u); res.Error != nil {
		return nil, res.Error
	}

	return u, nil
}
