package internal

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
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

	if issuer == "bookclub-access" { //TODO: constant
		expiresAt = now.Add(time.Duration(time.Hour) * 1)
	} else {
		expiresAt = now.Add(time.Hour * 1440)
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: issuer, IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(expiresAt), Subject: strconv.Itoa(id)})
	jwtSecret := os.Getenv("JWT_SECRET")
	signedJwt, err := jwt.SignedString([]byte(jwtSecret)) //TODO: change the jwt secret
	if err != nil {
		return "", errors.New("Coudn't create JWT token")
	}
	bearer := fmt.Sprintf("Bearer %s", signedJwt)
	return bearer, nil
}
