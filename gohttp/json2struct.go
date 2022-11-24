package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	data := []byte(`{
		"id": 2,
		"name": "david",
		"age": 35
	}`)

	//* ต้องส่งเป็น pointer เพราะ User อยู่ใน package main แต่ Unmarshal อยู่ใน package json ทำการส่ง reference ไปเพื่อให้ Unmarshal เอา data ใส่ใน u
	// var u User
	// json.Unmarshal(data, &u)
	// Or
	u := &User{}
	err := json.Unmarshal(data, u)

	fmt.Printf("% #v \n", u)
	fmt.Println(err)
}
