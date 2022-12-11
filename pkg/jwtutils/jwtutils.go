package jwtutils

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTParams struct {
	Scope     string
	Subject   string
	Issuer    string
	Secret    []byte
	ExpiresAt time.Time
}

type JWTInfo struct {
	Scope     string
	Subject   string
	Issuer    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func CreateJWT(params JWTParams) (string, error) {
	subject := fmt.Sprintf("%s:%s", params.Scope, params.Subject)
	claims := jwt.StandardClaims{
		ExpiresAt: params.ExpiresAt.Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    params.Issuer,
		Subject:   subject,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(params.Secret)
}

func ParseJWT(jwtString string, params JWTParams) (string, error) {
	token, err := jwt.ParseWithClaims(jwtString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return params.Secret, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("Invalid JWT")
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("Cannot obtain JWT claims")
	}
	i := strings.Index(claims.Subject, ":")
	scope := claims.Subject[:i]
	if params.Scope != scope {
		return "", fmt.Errorf(`JWT scope "%s" not as expected "%s"`, scope, params.Scope)
	}
	subject := claims.Subject[i+1:]
	return subject, nil
}

func ExtractInfoFromJWT(jwtString string) (*JWTInfo, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(jwtString, &jwt.StandardClaims{})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || claims.IssuedAt == 0 {
		return nil, fmt.Errorf("Cannot obtain JWT Info")
	}
	i := strings.Index(claims.Subject, ":")
	if i < 0 {
		return nil, fmt.Errorf("Invalid JWT Subject \"%s\"", claims.Subject)
	}
	scope := claims.Subject[:i]
	subject := claims.Subject[i:]
	return &JWTInfo{
		IssuedAt:  time.Unix(claims.IssuedAt, 0),
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
		Subject:   subject,
		Scope:     scope,
		Issuer:    claims.Issuer,
	}, nil
}
