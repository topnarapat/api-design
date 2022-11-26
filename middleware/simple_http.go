package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

func usersHandler(w http.ResponseWriter, r *http.Request) {
	// u, p, ok := r.BasicAuth()
	// log.Println("auth:", u, p, ok)

	if r.Method == "GET" {
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
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		next.ServeHTTP(w, r)
// 		log.Printf("Server http middleware: %s %s %s %s", r.RemoteAddr, r.Method, r.URL, time.Since(start))
// 	}
// }

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`can't parse the basic auth`))
			return
		}
		if u != "admin" || p != "1234" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`Username/Password incorrect`))
			return
		}
		fmt.Println("auth passed.")
		next(w, r)
	}
}

type Logger struct {
	Handler http.Handler
}

func (l Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.Handler.ServeHTTP(w, r)
	log.Printf("Server http middleware: %s %s %s %s", r.RemoteAddr, r.Method, r.URL, time.Since(start))
}

func main() {
	mux := http.NewServeMux()
	// http.HandleFunc("/users", logMiddleware(usersHandler))
	// http.HandleFunc("/health", logMiddleware(healthHandler))
	mux.HandleFunc("/users", AuthMiddleware(usersHandler))
	mux.HandleFunc("/health", healthHandler)

	logMux := Logger{Handler: mux}
	srv := http.Server{
		Addr:    ":2565",
		Handler: logMux,
	}

	log.Println("Server started at :2565")
	// log.Fatal(http.ListenAndServe(":2565", nil))
	// log.Fatal(http.ListenAndServe(":2565", mux))
	log.Fatal(srv.ListenAndServe())
	log.Println("bye bye!")
}
