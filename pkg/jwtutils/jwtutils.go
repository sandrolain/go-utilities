package jwtutils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTParams struct {
	Subject   string
	Issuer    string
	Secret    []byte
	ExpiresAt time.Time
}

type JWTInfo struct {
	Subject   string
	Issuer    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func CreateJWT(params JWTParams) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: params.ExpiresAt.Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    params.Issuer,
		Subject:   params.Subject,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(params.Secret)
}

func ParseJWT(jwtString string, params JWTParams) (string, error) {
	if jwtString == "" {
		return "", fmt.Errorf("the jwt string is empty")
	}
	token, err := jwt.ParseWithClaims(jwtString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return params.Secret, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid JWT")
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("cannot obtain JWT claims")
	}
	return claims.Subject, nil
}

func ExtractInfoFromJWT(jwtString string) (*JWTInfo, error) {
	if jwtString == "" {
		return nil, fmt.Errorf("the jwt string is empty")
	}
	token, _, err := new(jwt.Parser).ParseUnverified(jwtString, &jwt.StandardClaims{})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || claims.IssuedAt == 0 {
		return nil, fmt.Errorf("cannot obtain JWT Info")
	}
	return &JWTInfo{
		IssuedAt:  time.Unix(claims.IssuedAt, 0),
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
		Subject:   claims.Subject,
		Issuer:    claims.Issuer,
	}, nil
}
