package dao

import (
	"context"
	"time"
)

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

type SearchPatient struct {
	Firstname string
	Lastname  string
	Filters   string
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

type DBActions interface {
	AddStudy(ctx context.Context, study Study) error
	AddPatient(ctx context.Context, patient Patient) error
	GetPatient(ctx context.Context, filters SearchPatient, nestedValues bool) ([]Patient, error)
	GetPatientByUUID(ctx context.Context, uuid string, nestedValues bool) (Patient, error)
	GetStudyByUUID(ctx context.Context, uuid string) (Study, error)
	GetStudiesByPatientUUID(ctx context.Context, patientUUID string) ([]Study, error)
	DeletePatient(ctx context.Context, uuid string) error
	UpdatePatient(ctx context.Context, patient Patient) error
}
