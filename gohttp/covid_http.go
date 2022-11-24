package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type InfoCovid struct {
	ProvinceID int    `json:"ProvinceId"`
	Province   string `json:"Province"`
	Age        int    `json:"Age"`
}

type ResultCovid struct {
	Province map[string]int `json:"Province"`
	AgeGroup map[string]int `json:"AgeGroup"`
}

var x = []InfoCovid{
	{ProvinceID: 1, Province: "Bangkok", Age: 30},
	{ProvinceID: 1, Province: "Bangkok", Age: 30},
	{ProvinceID: 2, Province: "Chaing Mai", Age: 65},
	{ProvinceID: 2, Province: "Chaing Mai", Age: 50},
	{ProvinceID: 3, Province: "Phuket", Age: 20},
	{ProvinceID: 4, Province: "Krabi", Age: 0},
}

var y = ResultCovid{}

func main() {
	http.HandleFunc("/lmwn", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			a := ResultCovid{}
			m := make(map[string]int)
			n := make(map[string]int)
			for _, v := range x {
				m[v.Province] += 1
				if v.Age > 0 && v.Age <= 30 {
					n["0-30"] += 1
				}
				if v.Age >= 31 && v.Age <= 60 {
					n["31-60"] += 1
				}
				if v.Age >= 61 {
					n["61+"] += 1
				}
				if v.Age <= 0 {
					n["N/A"] += 1
				}
			}
			a.Province = m
			a.AgeGroup = n
			b, err := json.Marshal(a)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(b)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	log.Println("Server started at :2565")
	log.Fatal(http.ListenAndServe(":2565", nil))
	log.Println("bye bye!")
}
