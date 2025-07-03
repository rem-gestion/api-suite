// model/address.go
package model

import "time"

type Address struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Floor     string
	Unit      string
	Street    string
	Number    int
	City      string
	State     string
	Zip       string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
