package dao

import (
	"context"
	"time"
)

type Table struct {
	PK        string    `dynamodbav:"pk"`
	SK        string    `dynamodbav:"sk"`
	Fullname  string    `dynamodbav:"fullname"`
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
	Fullname  string    `dynamodbav:"fullname"`
	Filters   string    `dynamodbav:"filters"`
	CreatedAt time.Time `dynamodbav:"created_at"`
	UpdatedAt time.Time `dynamodbav:"updated_at"`
	DeletedAt time.Time `dynamodbav:"deleted_at"`
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
	GetPatientInfo(ctx context.Context, pk string) (*PatientInfo, error)
	GetPatientsInfo(ctx context.Context) ([]PatientInfo, error)
	SearchPatientInfo(ctx context.Context, fullname string) ([]PatientInfo, error)

	AddPatientInfo(ctx context.Context, data *PatientInfo) error

	UpdatePatientInfo(ctx context.Context, pk string, data *PatientInfo) error

	DeletePatient(ctx context.Context, pk string) error

	CheckDCM(ctx context.Context, hash, filename string) (bool, error)
	AddDCM(ctx context.Context, pk string, data *DCMInfo) error
}
