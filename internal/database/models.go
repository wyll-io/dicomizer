package database

import "time"

type patient struct {
	UUID      string    `dynamodbav:"uuid"`
	Firstname string    `dynamodbav:"firstname"`
	Lastname  string    `dynamodbav:"lastname"`
	Filters   string    `dynamodbav:"filters"`
	CreatedAt time.Time `dynamodbav:"created_at"`
	UpdatedAt time.Time `dynamodbav:"updated_at"`
	DeletedAt time.Time `dynamodbav:"deleted_at"`
}

type dcmImage struct {
	UUID        string    `dynamodbav:"uuid"`
	PatientUUID string    `dynamodbav:"patient_uuid"`
	Hash        string    `dynamodbav:"hash"`
	Filename    string    `dynamodbav:"filename"`
	CreatedAt   time.Time `dynamodbav:"created_at"`
	UpdatedAt   time.Time `dynamodbav:"updated_at"`
	DeletedAt   time.Time `dynamodbav:"deleted_at"`
}

type Dicom struct {
	UUID string `dynamodbav:"uuid"`
	UID  string `dynamodbav:"uid"`
}
