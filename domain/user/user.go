package user

type User struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
}

type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
}
