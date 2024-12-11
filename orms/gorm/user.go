package gorm

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	Id        string         `json:"id"`
	FirstName sql.NullString `json:"firstName" gorm:"column:first_name"`
	LastName  sql.NullString `json:"lastName" gorm:"column:last_name"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"-" gorm:"<-:create"`
	CreatedAt string         `json:"created_at" gorm:"column:created_at"`
}

type UserCreatePayload struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
}

type UserUpdatePayload struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
}

type IUserDao interface {
	List() ([]User, error)
	Get(string) (User, error)
	GetByEmail(string) (User, error)
	GetByUsername(string) (User, error)
	Create(UserCreatePayload) (User, error)
	Update(string, UserUpdatePayload) (User, error)
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
func (u *UserDao) Create(payload UserCreatePayload) (user User, err error) {
	user = User{
		FirstName: sql.NullString{String: *payload.FirstName, Valid: payload.FirstName != nil},
		LastName:  sql.NullString{String: *payload.LastName, Valid: payload.LastName != nil},
		Email:     payload.Email,
		Password:  payload.Password,
	}

	err = u.client.Create(&user).Error

	return user, err
}

func (u *UserDao) Update(id string, payload UserUpdatePayload) (user User, err error) {
	updates := make(map[string]interface{})
	if payload.FirstName != nil {
		updates["first_name"] = *payload.FirstName
	}
	if payload.LastName != nil {
		updates["last_name"] = *payload.LastName
	}
	if payload.Email != nil {
		updates["email"] = *payload.Email
	}
	if payload.Password != nil {
		updates["password"] = *payload.Password
	}

	err = u.client.Model(&User{}).Where("id = ?", id).Updates(updates).Error

	if err == nil {
		err = u.client.Where("id = ?", id).First(&user).Error
	}

	return user, err
}

func (u *UserDao) Delete(id string) error {
	return u.client.Where("id = ?", id).Delete(&User{}).Error
}

func NewUserDao(client *gorm.DB) IUserDao {
	return &UserDao{client: client}
}
