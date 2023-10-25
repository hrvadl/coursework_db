package repo

import (
	"github.com/hrvadl/coursework_db/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User interface {
	Get() ([]models.User, error)
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(u *models.User) (*models.User, error)
	Patch(u *models.User) (*models.User, error)
}

type stock struct {
	db *gorm.DB
}

func NewStock(db *gorm.DB) User {
	return &stock{db}
}

func (r *stock) Get() ([]models.User, error) {
	var users []models.User

	if tx := r.db.Where(&models.User{Role: models.StockRole}).Find(&users); tx.Error != nil {
		return nil, tx.Error
	}

	return users, nil
}

func (r *stock) GetByID(id int) (*models.User, error) {
	var user models.User
	res := r.db.Model(&models.User{}).
		Preload(clause.Associations).
		Preload("Inventory.Security").
		Where(&models.User{ID: uint(id), Role: models.StockRole}).
		First(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (r *stock) GetByEmail(email string) (*models.User, error) {
	var user models.User
	res := r.db.Where(&models.User{Email: email, Role: models.StockRole}).First(&user)

	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (r *stock) Create(u *models.User) (*models.User, error) {
	if res := r.db.Create(&u); res.Error != nil {
		return nil, res.Error
	}

	return u, nil
}

func (r *stock) Patch(u *models.User) (*models.User, error) {
	if res := r.db.Model(u).Updates(u); res.Error != nil {
		return nil, res.Error
	}

	return u, nil
}
