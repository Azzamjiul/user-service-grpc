// domain/user/repository/user_repository.go
package repository

import (
	"log"
	"user-service/domain/user"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) Create(user *user.User) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepository) FindByID(id uint) (*user.User, error) {
	var u user.User
	if err := ur.db.First(&u, id).Error; err != nil {
		log.Println("failed to query user: " + err.Error())
		return nil, err
	}

	log.Println("success to query user: ", u)
	return &u, nil
}

func (ur *UserRepository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	if err := ur.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
