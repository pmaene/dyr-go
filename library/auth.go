package library

import (
    jwt "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo"
)

type Scope int

const (
    AuthenticationScope = 0
    ActivationScope = 1
)

func GetTokenClaims(c echo.Context) jwt.MapClaims {
    user := c.Get("user").(*jwt.Token)
    return user.Claims.(jwt.MapClaims)
}
