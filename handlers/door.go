package handlers

import (
    "fmt"
    "net/http"

    "dyr/library"
    "dyr/models"

    "github.com/boltdb/bolt"
    "github.com/labstack/echo"
)

func GetDoors(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        doors, err := models.GetDoors(db)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

		return c.JSON(http.StatusOK, doors)
    }
}

func CreateDoor(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        claims := library.GetTokenClaims(c)
        if !claims["admin"].(bool) {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        bh := library.BoltHandler{db}
        bm := library.BoltModel{&models.Door{}}

        return bh.Create(&bm)(c)
    }
}
func GetDoor(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        claims := library.GetTokenClaims(c)
        if !claims["admin"].(bool) {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        bh := library.BoltHandler{db}
        bm := library.BoltModel{&models.Door{}}

        return bh.Get(&bm, "name")(c)
    }
}
func UpdateDoor(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        claims := library.GetTokenClaims(c)
        if !claims["admin"].(bool) {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        bh := library.BoltHandler{db}
        bm := library.BoltModel{&models.Door{}}

        return bh.Update(&bm, "name")(c)
    }
}
func DeleteDoor(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        claims := library.GetTokenClaims(c)
        if !claims["admin"].(bool) {
            return echo.NewHTTPError(http.StatusUnauthorized)
        }

        bh := library.BoltHandler{db}
        bm := library.BoltModel{&models.Door{}}

        return bh.Delete(&bm, "name")(c)
    }
}

func GetDoorNonce(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        var d models.Door
        if err := d.DecodeKey(c.Param("name")); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        if err := d.Get(db); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        return c.String(http.StatusOK, fmt.Sprintf("nonce/0x%08x", d.Nonce))
    }
}
func SwitchDoor(db *bolt.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        var d models.Door
        if err := d.DecodeKey(c.Param("name")); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        if err := d.Get(db); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        if err := d.Switch(db); err != nil {
            return echo.NewHTTPError(http.StatusServiceUnavailable)
        }

        claims := library.GetTokenClaims(c)

        u, err := models.GetTokenUser(db, claims)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        if err := models.CreateEvent(db, d.Name, u.Username); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        return c.NoContent(http.StatusOK)
    }
}
