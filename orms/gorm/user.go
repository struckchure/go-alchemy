package gorm

import "gorm.io/gorm"

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName" gorm:"column:first_name"`
	LastName  string `json:"lastName" gorm:"column:last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
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
	Id        *string `json:"id"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     *string `json:"email"`
	Username  *string `json:"username"`
	Password  *string `json:"password"`
}

type IUserDao interface {
	List() ([]User, error)
	Get(string) (User, error)
	GetByEmail(string) (User, error)
	GetByUsername(string) (User, error)
	Create(UserCreatePayload) (User, error)
	Update(UserUpdatePayload) (User, error)
	Delete(string) error
}

type UserDao struct {
	client *gorm.DB
}

func (u *UserDao) List() (users []User, err error) {
	err = u.client.Model(&User{}).Find(&users).Error

	return users, err
}

func (u *UserDao) Get(id string) (user User, err error) {
	err = u.client.Model(&User{}).Where("id = ?", id).First(&user).Error

	return user, err
}

func (u *UserDao) GetByEmail(email string) (user User, err error) {
	err = u.client.Model(&User{}).Where("email = ?", email).First(&user).Error

	return user, err
}

func (u *UserDao) GetByUsername(username string) (user User, err error) {
	err = u.client.Model(&User{}).Where("username = ?", username).First(&user).Error

	return user, err
}
func (u *UserDao) Create(payload UserCreatePayload) (user User, err error) {
	user = User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Username:  payload.Username,
		Password:  payload.Password,
	}

	err = u.client.Create(&user).Error

	return user, err
}

func (u *UserDao) Update(payload UserUpdatePayload) (user User, err error) {
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
	if payload.Username != nil {
		updates["username"] = *payload.Username
	}
	if payload.Password != nil {
		updates["password"] = *payload.Password
	}

	err = u.client.Model(&User{}).Where("id = ?", *payload.Id).Updates(updates).Error

	if err == nil {
		err = u.client.Where("id = ?", *payload.Id).First(&user).Error
	}

	return user, err
}

func (u *UserDao) Delete(id string) error {
	return u.client.Where("id = ?", id).Delete(&User{}).Error
}

func NewUserDao(client *gorm.DB) IUserDao {
	return &UserDao{client: client}
}
