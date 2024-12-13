// @alchemy replace package dao
package prisma

import (
	"context"
	"errors"

	// @alchemy statement "{{ .ModuleName }}/prisma/db"
	"github.com/struckchure/go-alchemy/prisma/db"
)

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"-"`
}

func (User) fromModel(user *db.UserModel) *User {
	if user == nil {
		return nil
	}

	firstName, _ := user.FirstName()
	lastName, _ := user.FirstName()

	return &User{
		Id:        user.ID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     user.Email,
		Password:  user.Password,
	}
}

type UserUpdatePayload struct {
	FirstName *string
	LastName  *string
	Email     *string
	Password  *string
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
	client *db.PrismaClient
}

func (u *UserDao) List() ([]User, error) {
	panic("unimplemented")
}

// @alchemy block {{- if .Login }}
func (u *UserDao) Get(id string) (*User, error) {
	panic("unimplemented")
}

func (u *UserDao) GetByEmail(email string) (*User, error) {
	ctx := context.Background()

	user, err := u.client.User.FindUnique(db.User.Email.Equals(email)).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return User{}.fromModel(user), err
}

// @alchemy block {{- end }}

// @alchemy block {{- if .Register }}
type UserCreatePayload struct {
	FirstName *string
	LastName  *string
	Email     string
	Password  string
}

func (u *UserDao) Create(payload UserCreatePayload) (*User, error) {
	ctx := context.Background()

	user, err := u.client.User.CreateOne(
		db.User.Email.Set(payload.Email),
		db.User.Password.Set(payload.Password),
		db.User.FirstName.SetIfPresent(payload.FirstName),
		db.User.LastName.SetIfPresent(payload.LastName),
	).Exec(ctx)

	if err != nil {
		if _, isUnique := db.IsErrUniqueConstraint(err); isUnique {
			return nil, errors.New("record already exist")
		}

		return nil, err
	}

	return User{}.fromModel(user), nil
}

// @alchemy block {{- end }}

func (u *UserDao) Update(id string, payload UserUpdatePayload) (*User, error) {
	panic("unimplemented")
}

func (u *UserDao) Delete(id string) error {
	panic("unimplemented")
}

func NewUserDao(client *db.PrismaClient) IUserDao {
	return &UserDao{client: client}
}
