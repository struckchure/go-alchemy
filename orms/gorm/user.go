// @alchemy replace package dao
package gorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	Id        string  `json:"id" gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	FirstName *string `json:"firstName" gorm:"column:first_name"`
	LastName  *string `json:"lastName" gorm:"column:last_name"`
	Email     string  `json:"email" gorm:"email;unique"`
	Password  string  `json:"-" gorm:"password"`
}

type IUserDao interface {
	List() ([]User, error)
	// @alchemy block {{- if .Login }}
	Get(string) (*User, error)
	GetByEmail(string) (*User, error)
	// @alchemy block {{- end }}
	// @alchemy block {{- if .Register }}
	Create(UserCreatePayload) (*User, error)
	// @alchemy block {{- end }}
	Update(string, UserUpdatePayload) (*User, error)
	Delete(string) error
}

type UserDao struct {
	client *gorm.DB
}

func (u *UserDao) List() (users []User, err error) {
	err = u.client.Model(&User{}).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, err
}

func (u *UserDao) Get(id string) (user *User, err error) {
	err = u.client.Model(&User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, err
}

// @alchemy block {{- if .Login }}
func (u *UserDao) GetByEmail(email string) (user *User, err error) {
	err = u.client.Model(&User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, err
}

// @alchemy block {{- end }}

// @alchemy block {{- if .Register }}

type UserCreatePayload struct {
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Email     string  `json:"email,omitempty"`
	Password  string  `json:"password,omitempty"`
}

func (u *UserDao) Create(payload UserCreatePayload) (*User, error) {
	user := User{}

	SetIfPresent(user, "FirstName", payload.FirstName)
	SetIfPresent(user, "LastName", payload.LastName)
	user.Email = payload.Email
	user.Password = payload.Password

	err := u.client.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}

// @alchemy block {{- end }}

type UserUpdatePayload struct {
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Email     *string `json:"email,omitempty"`
	Password  *string `json:"password,omitempty"`
}

func (u *UserDao) Update(id string, payload UserUpdatePayload) (*User, error) {
	user := User{}

	SetIfPresent(user, "FirstName", payload.FirstName)
	SetIfPresent(user, "LastName", payload.LastName)
	SetIfPresent(user, "Email", payload.Email)
	SetIfPresent(user, "Password", payload.Password)

	err := u.client.
		Clauses(clause.Returning{}).
		Model(&user).Where("id = ?", id).Updates(payload).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (u *UserDao) Delete(id string) error {
	return u.client.Where("id = ?", id).Delete(&User{}).Error
}

func NewUserDao(client *gorm.DB) IUserDao {
	return &UserDao{client: client}
}
