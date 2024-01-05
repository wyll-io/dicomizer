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
	Studies   []DCMImage
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type SearchPatientParams struct {
	Firstname string
	Lastname  string
	Filters   string
}

type DCMImage struct {
	UUID        string
	PatientUUID string
	Hash        string
	Filename    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type DBActions interface {
	AddStudy(ctx context.Context, study DCMImage) error
	AddPatient(ctx context.Context, patient Patient) error
	GetPatient(ctx context.Context, filters SearchPatientParams, nestedValues bool) ([]Patient, error)
	GetPatientByUUID(ctx context.Context, uuid string, nestedValues bool) (Patient, error)
	GetStudyByUUID(ctx context.Context, uuid string) (DCMImage, error)
	GetStudiesByPatientUUID(ctx context.Context, patientUUID string) ([]DCMImage, error)
	DeletePatient(ctx context.Context, uuid string) error
	UpdatePatient(ctx context.Context, patient Patient) error
}
