package handlers

import (
    "net/http"
    "time"

    "dyr/library"
    "dyr/models"

    "github.com/boltdb/bolt"
    "github.com/labstack/echo"
)

func Login(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        var p models.Password
        var u models.User

        if err := u.DecodeKey(c.FormValue("username")); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        if err := u.Get(db); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        if err := p.DecodeKey(u.Username); err != nil {
            return err
        }

        if err := p.Get(db); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        ok, err := p.Verify(c.FormValue("password"))
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        t, err := u.GetTokenString(
            time.Now().Add(time.Hour * 168),
            library.AuthenticationScope,
        )

        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        if ok {
    		return c.JSON(http.StatusOK, map[string]string{"token": t})
        }

        return echo.NewHTTPError(http.StatusUnauthorized)
    }
}
func Refresh(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        claims := library.GetTokenClaims(c)
        if library.Scope(claims["scope"].(float64)) != library.AuthenticationScope {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        u, err := models.GetTokenUser(db, claims)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        t, err := u.GetTokenString(
            time.Now().Add(time.Hour * 168),
            library.AuthenticationScope,
        )

        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

    	return c.JSON(http.StatusOK, map[string]string{"token": t})
    }
}
