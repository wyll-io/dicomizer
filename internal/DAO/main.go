package dao

import (
	"context"
	"time"
)

type Table struct {
	PK        string    `dynamodbav:"pk"`
	SK        string    `dynamodbav:"sk"`
	Firstname string    `dynamodbav:"firstname"`
	Lastname  string    `dynamodbav:"lastname"`
	Filters   string    `dynamodbav:"filters"`
	Hash      string    `dynamodbav:"hash"`
	Filename  string    `dynamodbav:"filename"`
	CreatedAt time.Time `dynamodbav:"created_at"`
	UpdatedAt time.Time `dynamodbav:"updated_at"`
	DeletedAt time.Time `dynamodbav:"deleted_at"`
}

type PatientInfo struct {
	PK        string    `dynamodbav:"pk"`
	SK        string    `dynamodbav:"sk"`
	Firstname string    `dynamodbav:"firstname"`
	Lastname  string    `dynamodbav:"lastname"`
	Filters   string    `dynamodbav:"filters"`
	CreatedAt time.Time `dynamodbav:"created_at"`
	UpdatedAt time.Time `dynamodbav:"updated_at"`
	DeletedAt time.Time `dynamodbav:"deleted_at"`

	DCMCount uint
}

type DCMInfo struct {
	PK        string    `dynamodbav:"pk"`
	SK        string    `dynamodbav:"sk"`
	Hash      string    `dynamodbav:"hash"`
	Filename  string    `dynamodbav:"filename"`
	CreatedAt time.Time `dynamodbav:"created_at"`
	DeletedAt time.Time `dynamodbav:"deleted_at"`
}

type DBActions interface {
	AddPatientInfo(ctx context.Context, data *PatientInfo) error
	AddPatientDCM(ctx context.Context, pk string, data *DCMInfo) error
	SearchPatientInfo(ctx context.Context, fullname string) ([]PatientInfo, error)
	GetPatientsInfo(ctx context.Context) ([]PatientInfo, error)
	UpdatePatientInfo(ctx context.Context, pk string, data *PatientInfo) error
	DeletePatient(ctx context.Context, pk string) error
}
