package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJwtService interface {
	GenerateAccessToken(Claims) (*string, error)
	GenerateRefreshToken(Claims) (*string, error)
	GenerateTokens(Claims) (*Tokens, error)
	ValidateAccessToken(string) (*Claims, error)
	ValidateRefreshToken(string) (*Claims, error)
}

type JwtService struct{}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Claims struct {
	Sub                  string `json:"sub"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // Use this for standard fields like exp, iss, etc.
}

var (
	JWT_ACCESS_TOKEN_SECRET  string = GetEnv("JWT_ACCESS_TOKEN_SECRET", "access-token")
	JWT_ACCESS_TOKEN_EXPIRY  string = GetEnv("JWT_ACCESS_TOKEN_EXPIRY", "10m")
	JWT_REFRESH_TOKEN_SECRET string = GetEnv("JWT_ACCESS_TOKEN_SECRET", "refresh-secret")
	JWT_REFRESH_TOKEN_EXPIRY string = GetEnv("JWT_ACCESS_TOKEN_EXPIRY", "168h") // 7d
)

func (j *JwtService) generateToken(claims Claims, secret string, expiry string) (*string, error) {
	parsedExpiry, err := time.ParseDuration(expiry)
	if err != nil {
		return nil, err
	}
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(parsedExpiry))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(secret))

	return &signedToken, nil
}

func (j *JwtService) validateToken(token string, secret string) (*Claims, error) {
	validatedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := validatedToken.Claims.(*Claims)
	if !ok && !validatedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (j *JwtService) GenerateAccessToken(claims Claims) (*string, error) {
	return j.generateToken(claims, JWT_ACCESS_TOKEN_SECRET, JWT_ACCESS_TOKEN_EXPIRY)
}

func (j *JwtService) GenerateRefreshToken(claims Claims) (*string, error) {
	return j.generateToken(claims, JWT_REFRESH_TOKEN_SECRET, JWT_REFRESH_TOKEN_EXPIRY)
}

func (j *JwtService) GenerateTokens(claims Claims) (*Tokens, error) {
	accessToken, err := j.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.GenerateRefreshToken(claims)
	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: *accessToken, RefreshToken: *refreshToken}, nil
}

func (j *JwtService) ValidateAccessToken(token string) (*Claims, error) {
	return j.validateToken(token, JWT_ACCESS_TOKEN_SECRET)
}

func (j *JwtService) ValidateRefreshToken(token string) (*Claims, error) {
	return j.validateToken(token, JWT_REFRESH_TOKEN_SECRET)
}

func NewJwtService() IJwtService {
	return &JwtService{}
}
