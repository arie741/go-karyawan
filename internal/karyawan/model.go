package karyawan

import "time"

type Karyawan struct {
	Id        string    `bson:"Id"`
	Name      string    `bson:"Name"`
	BirthDate time.Time `bson:"BirthDate"`
	Salary    int       `bson:"Salary"`
	Position  string    `bson:"Position"`
	Joined    time.Time `bson:"Joined"`
}
