package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
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
	stmt, err := db.Prepare("SELECT id, name, age FROM users")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "can't prepare query all users statement"})
	}
	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "can't query all users:" + err.Error()})
	}
	users := []User{}
	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Name, &u.Age)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: "can't scan users:" + err.Error()})
		}
		users = append(users, u)
	}

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

	row := db.QueryRow("INSERT INTO users (name, age) values ($1, $2) RETURNING id", u.Name, u.Age)
	err = row.Scan(&u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	// users = append(users, u)

	return c.JSON(http.StatusCreated, u)
}

func getUserHandler(c echo.Context) error {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, name, age FROM users WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query user statment:" + err.Error()})
	}

	row := stmt.QueryRow(id)
	u := User{}
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

func updateUserHandler(c echo.Context) error {
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

func deleteUserHandler(c echo.Context) error {
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

var db *sql.DB

func main() {
	var err error
	godotenv.Load("../.env")
	db, err = sql.Open("postgres", os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	createTb := `CREATE TABLE IF NOT EXISTS users ( id SERIAL PRIMARY KEY, name TEXT, age INT)`

	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", healthHandler)

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "1234" {
			return true, nil
		}
		return false, nil
	}))

	e.GET("/users", getUsersHandler)
	e.GET("/users/:id", getUserHandler)
	e.POST("/users", createUsersHandler)
	e.PATCH("/users/:id", updateUserHandler)
	e.DELETE("/users/:id", deleteUserHandler)

	log.Println("Server started at :2565")
	log.Fatal(e.Start(":2565"))
	log.Println("bye bye!")
}
