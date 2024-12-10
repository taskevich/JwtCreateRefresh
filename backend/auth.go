package backend

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	SECRETKEY        = "SUPER-SECRET"
	ACCESSTOKENTIME  = time.Hour * 4
	REFRESHTOKENTIME = time.Hour
)

func CreateToken(user *UserItem, tokenExpireTime time.Duration, guid *string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user": user,
		"guid": guid,
		"exp":  time.Now().Add(tokenExpireTime).Unix(),
		"iat":  time.Now().Unix(),
	})

	token, err := claims.SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRETKEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return token, nil
}

func GenerateAuthRefreshTokens(user *UserItem) (string, string, error) {
	guid := uuid.New().String()
	token, err := CreateToken(user, ACCESSTOKENTIME, &guid)
	if err != nil {
		return "", "", err
	}

	refresh, err := CreateToken(user, REFRESHTOKENTIME, &guid)
	if err != nil {
		return "", "", err
	}
	return token, refresh, nil
}
