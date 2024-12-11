package services

import (
	"errors"

	// @alchemy statement "{{ .ModuleName }}/dao"
	"github.com/struckchure/go-alchemy/orms/prisma"
	// @alchemy statement "{{ .ModuleName }}/prisma/db"
	"github.com/struckchure/go-alchemy/prisma/db"
)

type IAuthenticationService interface {
	// @alchemy block {{- if .Login }}
	Login(LoginArgs) (*LoginResult, error)
	// @alchemy block {{- end }}
	// @alchemy block {{- if .Register }}
	Register(RegisterArgs) (*RegisterResult, error)
	// @alchemy block {{- end }}
}

type AuthenticationService struct {
	// @alchemy replace userDao dao.IUserDao
	userDao    prisma.IUserDao
	jwtService IJwtService
}

func (a *AuthenticationService) passwordIsValid(plain string, hashed string) (bool, error) {
	return plain == hashed, nil
}

// @alchemy block {{- if .Login  }}
type LoginArgs struct {
	Email    string
	Password string
}

type LoginResult struct {
	// @alchemy replace User dao.User `json:"user"`
	User   prisma.User `json:"user"`
	Tokens Tokens      `json:"tokens"`
}

func (a *AuthenticationService) Login(args LoginArgs) (*LoginResult, error) {
	user, err := a.userDao.GetByEmail(args.Email)
	if err != nil {
		if errors.Is(db.ErrNotFound, err) {
			return nil, errors.New("invalid credentials")
		}

		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	passwordIsValid, err := a.passwordIsValid(args.Password, user.Password)
	if err != nil {
		return nil, err
	}

	if !passwordIsValid {
		return nil, errors.New("invalid credentials")
	}

	tokens, err := a.jwtService.GenerateTokens(Claims{Sub: user.Id})
	if err != nil {
		return nil, err
	}

	return &LoginResult{User: *user, Tokens: *tokens}, nil
}

// @alchemy block {{- end }}

type RegisterArgs struct {
	FirstName *string
	LastName  *string
	Email     string
	Password  string
}

// @alchemy block {{- if .Register }}

type RegisterResult struct {
	// @alchemy replace User dao.User `json:"user"`
	User   prisma.User `json:"user"`
	Tokens Tokens      `json:"tokens"`
}

func (a *AuthenticationService) Register(args RegisterArgs) (*RegisterResult, error) {
	user, err := a.userDao.Create(
		// @alchemy replace dao.UserCreatePayload{
		prisma.UserCreatePayload{
			FirstName: args.FirstName,
			LastName:  args.LastName,
			Email:     args.Email,
			Password:  args.Password,
		},
	)
	if err != nil {
		if _, isUniqueConstraint := db.IsErrUniqueConstraint(err); isUniqueConstraint {
			return nil, errors.New("user already exist")
		}

		return nil, err
	}

	tokens, err := a.jwtService.GenerateTokens(Claims{Sub: user.Id})
	if err != nil {
		return nil, err
	}

	return &RegisterResult{User: *user, Tokens: *tokens}, nil
}

// @alchemy block {{- end }}

func NewAuthenticationService(
	// @alchemy replace userDao dao.IUserDao,
	userDao prisma.IUserDao,
	jwtService IJwtService,
) IAuthenticationService {
	return &AuthenticationService{
		userDao:    userDao,
		jwtService: jwtService,
	}
}
