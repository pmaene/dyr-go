package main

import (
    "dyr/handlers"

    "github.com/boltdb/bolt"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
)

const (
    apiRoute = "/api/v1"
)

func main() {
	e := echo.New()

    // Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

    e.Pre(middleware.RemoveTrailingSlash())

    // Database
    db, err := bolt.Open("dyr.db", 0600, nil)
    if err != nil {
        e.Logger.Fatal(err)
    }
    defer db.Close()

    // Routes
	r := e.Group("")
	r.Use(middleware.JWT([]byte("secret")))

    e.POST(apiRoute + "/auth/login", handlers.Login(db))
    r.GET(apiRoute + "/auth/refresh", handlers.Refresh(db))

    r.GET(apiRoute + "/doors", handlers.GetDoors(db))
    r.POST(apiRoute + "/doors", handlers.CreateDoor(db))
    r.GET(apiRoute + "/doors/:door", handlers.GetDoor(db))
    r.PUT(apiRoute + "/doors/:door", handlers.UpdateDoor(db))
    r.DELETE(apiRoute + "/doors/:door", handlers.DeleteDoor(db))

    e.GET(apiRoute + "/doors/:name/nonce", handlers.GetDoorNonce(db))
    r.POST(apiRoute + "/doors/:name/switch", handlers.SwitchDoor(db))

    r.GET(apiRoute + "/users", handlers.GetUsers(db))
    r.POST(apiRoute + "/users", handlers.CreateUser(db))
    r.GET(apiRoute + "/users/:username", handlers.GetUser(db))
    r.PUT(apiRoute + "/users/:username", handlers.UpdateUser(db))
    r.DELETE(apiRoute + "/users/:username", handlers.DeleteUser(db))

    r.GET(apiRoute + "/users/:username/activationToken", handlers.GetActivationToken(db))
    r.PUT(apiRoute + "/users/activate", handlers.ActivateUser(db))

	// Go
	e.Logger.Fatal(e.Start(":1323"))
}
