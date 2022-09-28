package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func main() {

	db, err := gorm.Open(sqlite.Open("tests.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	e.Use(middleware.Logger())
	e.Static("/", "./client")
	e.GET("/context-cancel", func(c echo.Context) (err error) {
		time.Sleep(2 * time.Second)
		log.Print("requested")
		tx := db.WithContext(c.Request().Context())
		user := &User{
			Name:  "test",
			Email: "test",
		}
		tx = tx.Create(&user)
		log.Print(tx.Error)
		return c.JSON(200, user)
	})
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			done := make(chan bool)
			go func() {
				next(c)
				done <- true
			}()
			ctx := c.Request().Context()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-done:
				log.Print("done")
			}

			return ctx.Err()
		}
	})
	e.Start(":3000")
}
