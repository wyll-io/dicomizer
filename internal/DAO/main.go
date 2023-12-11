package dao

import "time"

type Patient struct {
	UUID      string
	Firstname string
	Lastname  string
	Filters   string
	Studies   []Study
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Study struct {
	UUID        string
	PatientUUID string
	Status      string
	Hash        string
	Filename    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
