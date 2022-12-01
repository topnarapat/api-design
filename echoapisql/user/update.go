package user

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateUserHandler(c echo.Context) error {
	id := c.Param("id")

	var u User
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("UPDATE users SET name=$2, age=$3 WHERE id=$1 RETURNING id, name, age", id, u.Name, u.Age)
	err = row.Scan(&u.ID, &u.Name, &u.Age)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "user not found"})
	case nil:
		return c.JSON(http.StatusOK, u)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan user:" + err.Error()})
	}
}
