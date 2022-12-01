//go:build integration

package user

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGetAllUser(t *testing.T) {
	seedUser(t)

	var us []User
	res := request(http.MethodGet, uri("users"), nil)
	err := res.Decode(&us)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(us), 0)
}

func TestCreateUser(t *testing.T) {
	body := bytes.NewBufferString(`{
		"name": "anuchito",
		"age": 19
	}`)

	var u User
	res := request(http.MethodPost, uri("users"), body)
	err := res.Decode(&u)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, u.ID)
	assert.Equal(t, "anuchito", u.Name)
	assert.Equal(t, 19, u.Age)
}

func TestGetUserByID(t *testing.T) {
	c := seedUser(t)

	var latest User
	res := request(http.MethodGet, uri("users", strconv.Itoa(c.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, latest.ID)
	assert.NotEmpty(t, latest.Name)
	assert.NotEmpty(t, latest.Age)
}

func TestUpdateUserByID(t *testing.T) {
	// t.Skip("TODO: implement me")
	c := seedUser(t)
	body := bytes.NewBufferString(`{
		"name": "Top",
		"age": 30
	}`)

	var u User
	r := request(http.MethodPatch, uri("users", strconv.Itoa(c.ID)), body)
	err := r.Decode(&u)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, c.ID, u.ID)
	assert.Equal(t, "Top", u.Name)
	assert.Equal(t, 30, u.Age)
}

func TestDeleteUserByID(t *testing.T) {
	// t.Skip("TODO: implement me")
	c := seedUser(t)

	var latest User
	res := request(http.MethodDelete, uri("users", strconv.Itoa(c.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Empty(t, latest.ID)
	assert.Empty(t, latest.Name)
	assert.Empty(t, latest.Age)
}

func seedUser(t *testing.T) User {
	var c User
	body := bytes.NewBufferString(`{
		"name": "Tangfa",
		"age": 29
	}`)
	err := request(http.MethodPost, uri("users"), body).Decode(&c)
	if err != nil {
		t.Fatal("can't create uomer:", err)
	}

	return c
}

func uri(path ...string) string {
	host := "http://localhost:2565"
	if path == nil {
		return host
	}

	url := append([]string{host}, path...)
	return strings.Join(url, "/")
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func request(method, url string, body io.Reader) *Response {
	godotenv.Load("../../.env")
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
