package dto

import "time"

type RequestCreateUser struct {
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	BirthDate  time.Time `json:"birth_date"`
	Password   string    `json:"password"`
}
