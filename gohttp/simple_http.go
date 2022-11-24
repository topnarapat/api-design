package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			log.Println("GET")
			//* covert string to byte
			// w.Write([]byte(`{"name": "fatcat", "method": "GET"}`))
			b, err := json.Marshal(users)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(b)
			return
		}

		if r.Method == "POST" {
			log.Println("POST")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			var u User
			err = json.Unmarshal(body, &u)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			users = append(users, u)

			fmt.Fprintf(w, "hello %s created users", "POST")
			// w.Write([]byte(`{"name": "chopper", "method": "POST"}`))
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	log.Println("Server started at :2565")
	log.Fatal(http.ListenAndServe(":2565", nil))
	log.Println("bye bye!")
}
