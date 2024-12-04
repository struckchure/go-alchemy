package services

import (
	"github.com/struckchure/go-alchemy"
	"github.com/struckchure/go-alchemy/orms/prisma"
	"github.com/struckchure/go-alchemy/prisma/db"
)

type IAuthenticationService interface {
	Login(LoginArgs) (*LoginResult, error)
	Register(RegisterArgs) (*RegisterResult, error)
}

type AuthenticationService struct {
	userDao prisma.IUserDao
}

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginArgs struct {
	Email    string
	Password string
}

type LoginResult struct {
	User   prisma.User `json:"user"`
	Tokens Token       `json:"tokens"`
}

func (a *AuthenticationService) passwordIsValid(plain string, hashed string) (bool, error) {
	return plain == hashed, nil
}

func (a *AuthenticationService) Login(args LoginArgs) (*LoginResult, error) {
	user, err := a.userDao.GetByEmail(args.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, alchemy.ErrInvalidCredentials
	}

	passwordIsValid, err := a.passwordIsValid(args.Password, user.Password)
	if err != nil {
		return nil, err
	}

	if !passwordIsValid {
		return nil, alchemy.ErrInvalidCredentials
	}

	tokens := Token{
		AccessToken:  "exampleAccessToken",
		RefreshToken: "exampleRefreshToken",
	}

	return &LoginResult{
		User:   *user,
		Tokens: tokens,
	}, nil
}

type RegisterArgs struct {
	FirstName *string
	LastName  *string
	Email     string
	Password  string
}

type RegisterResult struct {
	User   prisma.User `json:"user"`
	Tokens Token       `json:"tokens"`
}

func (a *AuthenticationService) Register(args RegisterArgs) (*RegisterResult, error) {
	user, err := a.userDao.Create(prisma.UserCreatePayload{
		FirstName: args.FirstName,
		LastName:  args.LastName,
		Email:     args.Email,
		Password:  args.Password,
	})
	if err != nil {
		return nil, err
	}

	tokens := Token{
		AccessToken:  "",
		RefreshToken: "",
	}

	return &RegisterResult{
		User:   *user,
		Tokens: tokens,
	}, nil
}

func NewAuthenticationService(client *db.PrismaClient) IAuthenticationService {
	return &AuthenticationService{
		userDao: prisma.NewUserDao(client),
	}
}
