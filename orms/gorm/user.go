package gorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName" gorm:"column:first_name"`
	LastName  string `json:"lastName" gorm:"column:last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at" gorm:"column:created_at"`
}

type UserCreatePayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type UserUpdatePayload struct {
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Email     *string `json:"email,omitempty"`
	Username  *string `json:"username,omitempty"`
	Password  *string `json:"password,omitempty"`
}

type IUserDao interface {
	List() ([]User, error)
	Get(string) (*User, error)
	GetByEmail(string) (*User, error)
	GetByUsername(string) (*User, error)
	Create(UserCreatePayload) (*User, error)
	Update(string, UserUpdatePayload) (*User, error)
	Delete(string) error
}

type UserDao struct {
	client *gorm.DB
}

func (u *UserDao) List() (users []User, err error) {
	err = u.client.Model(&User{}).Find(&users).Error

	return users, err
}

func (u *UserDao) Get(id string) (user *User, err error) {
	err = u.client.Model(&User{}).Where("id = ?", id).First(&user).Error

	return user, err
}

func (u *UserDao) GetByEmail(email string) (user *User, err error) {
	err = u.client.Model(&User{}).Where("email = ?", email).First(&user).Error

	return user, err
}

func (u *UserDao) GetByUsername(username string) (user *User, err error) {
	err = u.client.Model(&User{}).Where("username = ?", username).First(&user).Error

	return user, err
}

func (u *UserDao) Create(payload UserCreatePayload) (user *User, err error) {
	user = &User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Username:  payload.Username,
		Password:  payload.Password,
	}

	err = u.client.Create(&user).Error

	return user, err
}

func (u *UserDao) Update(id string, payload UserUpdatePayload) (user *User, err error) {
	err = u.client.
		Clauses(clause.Returning{}).
		Model(&user).
		Where("id = ?", id).
		Updates(payload).Error

	return user, err
}

func (u *UserDao) Delete(id string) error {
	return u.client.Where("id = ?", id).Delete(&User{}).Error
}

func NewUserDao(client *gorm.DB) IUserDao {
	return &UserDao{client: client}
}
