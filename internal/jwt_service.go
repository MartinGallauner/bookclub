package internal

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type JwtService interface {
	CreateToken(issuer string, id int) (string, error)
}

type JwtServiceImpl struct {
}

func (srv JwtServiceImpl) CreateToken(issuer string, id int) (string, error) {

	now := time.Now()
	var expiresAt time.Time

	if issuer == "bookclub-access" { //todo constant
		expiresAt = now.Add(time.Duration(time.Hour) * 1)
	} else {
		expiresAt = now.Add(time.Hour * 1440)
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: issuer, IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(expiresAt), Subject: strconv.Itoa(id)})
	signedJwt, err := jwt.SignedString([]byte("CHANGE_ME")) //todo change the jwt secret

	if err != nil {
		return "", errors.New("Coudn't create JWT token")
	}

	return signedJwt, nil
}
