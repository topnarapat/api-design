package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	u := User{
		ID: 1, Name: "john", Age: 30,
	}
	//* convert struct to json
	b, err := json.Marshal(u)
	fmt.Printf("byte : %T \n", b)
	fmt.Printf("byte : %v \n", b)
	fmt.Printf("byte : %s \n", b)
	fmt.Println(err)
}
