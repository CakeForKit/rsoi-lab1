package model

type Person struct {
	BaseEntity
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string
	Work    string
}
