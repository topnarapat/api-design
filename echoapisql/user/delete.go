package user

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func DeleteUserHandler(c echo.Context) error {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROm users WHERE id=$1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query delete user statment:" + err.Error()})
	}

	var u User
	row := stmt.QueryRow(id)
	err = row.Scan(&u.ID, &u.Name, &u.Age)

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "user not found"})
	case nil:
		return c.JSON(http.StatusOK, u)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan user:" + err.Error()})
	}
}
