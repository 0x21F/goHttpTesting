package models

type Person struct {
	Id    uint   `json:"id"`
	First string `json:"first"`
	Last  string `json:"last"`
}
