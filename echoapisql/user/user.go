package user

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Err struct {
	Message string `json:"message"`
}

// var users = []User{
// 	{ID: 1, Name: "john", Age: 30},
// 	{ID: 2, Name: "david", Age: 35},
// }
