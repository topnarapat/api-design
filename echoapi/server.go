package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{
	{ID: 1, Name: "john", Age: 30},
	{ID: 2, Name: "david", Age: 35},
}

func healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func getUsersHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

type Err struct {
	Message string `json:"message"`
}

func createUsersHandler(c echo.Context) error {
	var u User
	//* ส่ง address ของ u ไปเพราะต้องการ assign ค่า u และ u กับ func Bind อยู่คนละ package กัน
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	users = append(users, u)

	return c.JSON(http.StatusCreated, u)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", healthHandler)

	g := e.Group("/api")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "1234" {
			return true, nil
		}
		return false, nil
	}))

	g.GET("/users", getUsersHandler)
	g.POST("/users", createUsersHandler)

	log.Println("Server started at :2565")
	log.Fatal(e.Start(":2565"))
	log.Println("bye bye!")
}
