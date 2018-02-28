package handlers

import (
    "net/http"
    "time"

    "dyr/library"
    "dyr/models"

    "github.com/boltdb/bolt"
    "github.com/labstack/echo"
)

func GetUsers(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
	    claims := library.GetTokenClaims(c)
        if !claims["admin"].(bool) {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        users, err := models.GetUsers(db)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

		return c.JSON(http.StatusOK, users)
    }
}

func CreateUser(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        claims := library.GetTokenClaims(c)
        if !claims["admin"].(bool) {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        bh := library.BoltHandler{db}
        bm := library.BoltModel{&models.User{}}

        return bh.Create(&bm)(c)
    }
}
func GetUser(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        claims := library.GetTokenClaims(c)
        if !claims["admin"].(bool) {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        bh := library.BoltHandler{db}
        bm := library.BoltModel{&models.User{}}

        return bh.Get(&bm, "username")(c)
    }
}
func UpdateUser(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        claims := library.GetTokenClaims(c)
        if !claims["admin"].(bool) {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        bh := library.BoltHandler{db}
        bm := library.BoltModel{&models.User{}}

        return bh.Update(&bm, "username")(c)
    }
}
func DeleteUser(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        claims := library.GetTokenClaims(c)
        if !claims["admin"].(bool) {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        bh := library.BoltHandler{db}
        bm := library.BoltModel{&models.User{}}

        return bh.Delete(&bm, "username")(c)
    }
}

func GetActivationToken(db *bolt.DB) echo.HandlerFunc {
    return func (c echo.Context) error {
        claims := library.GetTokenClaims(c)
        if !claims["admin"].(bool) {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        var u models.User
        u.DecodeKey(c.Param("username"))
        if err := u.Get(db); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        t, err := u.GetTokenString(
            time.Now().Add(time.Hour * 24),
            library.ActivationScope,
        )

        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        return c.JSON(http.StatusOK, map[string]string{"token": t})
    }
}
func ActivateUser(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        claims := library.GetTokenClaims(c)
        if library.Scope(claims["scope"].(float64)) != library.ActivationScope {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        u := claims["sub"].(string)

        p, err := models.NewPassword(db, u, c.FormValue("password"))
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        if err := p.Put(db); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        return c.NoContent(http.StatusOK)
    }
}
