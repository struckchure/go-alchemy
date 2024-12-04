package gorm

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type UserCreatePayload struct{}

type UserUpdatePayload struct{}

type IUserDao interface {
	List() ([]User, error)
	Get(string) (*User, error)
	GetByEmail(string) (*User, error)
	GetByUsername(string) (*User, error)
	Create(UserCreatePayload) (*User, error)
	Update(UserUpdatePayload) (*User, error)
	Delete(string) error
}

type UserDao struct{}

func (u *UserDao) List() ([]User, error) {
	panic("unimplemented")
}

func (u *UserDao) Get(id string) (*User, error) {
	panic("unimplemented")
}

func (u *UserDao) GetByEmail(email string) (*User, error) {
	panic("unimplemented")
}

func (u *UserDao) GetByUsername(username string) (*User, error) {
	panic("unimplemented")
}

func (u *UserDao) Create(payload UserCreatePayload) (*User, error) {
	panic("unimplemented")
}

func (u *UserDao) Update(payload UserUpdatePayload) (*User, error) {
	panic("unimplemented")
}

func (u *UserDao) Delete(id string) error {
	panic("unimplemented")
}

func NewUserDao() IUserDao {
	return &UserDao{}
}
