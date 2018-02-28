package library

import (
    "net/http"

    "github.com/boltdb/bolt"
    "github.com/labstack/echo"
)

type BoltHandler struct {
    DB *bolt.DB
}
func (bh BoltHandler) Create(bm *BoltModel) echo.HandlerFunc {
    return func(c echo.Context) error {
        if err := c.Bind(&bm.Value); err != nil {
            return err
        }

        if err := bm.Put(bh.DB); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        return c.JSON(http.StatusCreated, bm.Value)
    }
}
func (bh BoltHandler) Get(bm *BoltModel, k string) echo.HandlerFunc {
    return func(c echo.Context) error {
        if err := get(bh.DB, bm, c.Param(k)); err != nil {
            if err == BucketDoesNotExistError || err == KeyNotFoundError {
                return echo.NewHTTPError(http.StatusNotFound)
            }

            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        return c.JSON(http.StatusOK, bm.Value)
    }
}
func (bh BoltHandler) Update(bm *BoltModel, k string) echo.HandlerFunc {
    return func(c echo.Context) error {
        if err := get(bh.DB, bm, c.Param(k)); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        if err := c.Bind(&bm.Value); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        if err := bm.Put(bh.DB); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        return c.JSON(http.StatusOK, bm.Value)
    }
}
func (bh BoltHandler) Delete(bm *BoltModel, k string) echo.HandlerFunc {
    return func(c echo.Context) error {
        if err := get(bh.DB, bm, c.Param(k)); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        if err := bm.Delete(bh.DB); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError)
        }

        return c.NoContent(http.StatusNoContent)
    }
}

func get(db *bolt.DB, bm *BoltModel, k string) error {
    if err := bm.Value.DecodeKey(k); err != nil {
        return err
    }

    if err := bm.Get(db); err != nil {
        return err
    }

    return nil
}
