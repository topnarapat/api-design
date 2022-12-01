package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateUsersHandler(c echo.Context) error {
	var u User
	//* ส่ง address ของ u ไปเพราะต้องการ assign ค่า u และ u กับ func Bind อยู่คนละ package กัน
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("INSERT INTO users (name, age) values ($1, $2) RETURNING id", u.Name, u.Age)
	err = row.Scan(&u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	// users = append(users, u)

	return c.JSON(http.StatusCreated, u)
}
