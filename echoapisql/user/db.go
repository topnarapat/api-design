package user

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	var err error
	godotenv.Load("../.env")
	db, err = sql.Open("postgres", os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTb := `CREATE TABLE IF NOT EXISTS users ( id SERIAL PRIMARY KEY, name TEXT, age INT)`

	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
}
