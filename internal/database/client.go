package database

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	dao "github.com/wyll-io/dicomizer/internal/DAO"
)

type DB struct {
	Client *dynamodb.Client
	Table  string
}

// New returns a new instance of the DynamoDB database client. It requires a
// valid AWS configuration and the name of the table to use.
// It uses a single table design, where the partition key is the UUID of the
// patient and the sort key is the UUID of the DICOM file.
func New(cfg aws.Config, table string) dao.DBActions {
	return DB{
		Client: dynamodb.NewFromConfig(cfg),
		Table:  table,
	}
}

// AddPatient adds a new patient to the database.
// PK, SK, CreatedAt, UpdatedAt and DeletedAt are automatically populated.
func (db DB) AddPatientInfo(ctx context.Context, data *dao.PatientInfo) error {
	data.PK = fmt.Sprintf("PATIENT#%s", uuid.New())
	data.SK = "INFO#0"
	data.CreatedAt = time.Now()

	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		return err
	}

	_, err = db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &db.Table,
		Item:      item,
	})

	return err
}

// AddPatientDCM adds a new DICOM file to the database.
// PK, SK, CreatedAt and DeletedAt are automatically populated.
func (db DB) AddPatientDCM(ctx context.Context, pk string, data *dao.DCMInfo) error {
	data.PK = pk
	data.SK = fmt.Sprintf("DCM#%s", uuid.New())
	data.CreatedAt = time.Now()

	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		return err
	}

	_, err = db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &db.Table,
		Item:      item,
	})

	return err
}

// SearchPatientInfo searches for patient info by fullname
// (case sensitive, dynamodb doesn't implement full-text search).
func (db DB) SearchPatientInfo(ctx context.Context, fullname string) ([]dao.PatientInfo, error) {
	// search for patient info
	pkEx := expression.Key("pk").BeginsWith("PATIENT#")
	skEx := expression.Key("sk").Equal(expression.Value("INFO#0"))
	fullnameEx := expression.Key("fullname").BeginsWith(fullname)

	expr, err := expression.NewBuilder().
		WithKeyCondition(pkEx.And(skEx).And(fullnameEx)).
		Build()
	if err != nil {
		return nil, err
	}

	res, err := db.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &db.Table,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return []dao.PatientInfo{}, nil
	}

	patients := make([]dao.PatientInfo, 0, res.Count)
	for _, i := range res.Items {
		patient := dao.PatientInfo{}
		if err := attributevalue.UnmarshalMap(i, &patient); err != nil {
			return nil, err
		}

		count, err := db.countPatientDCM(ctx, patient.PK)
		if err != nil {
			return nil, err
		}

		patient.DCMCount = count

		patients = append(patients, patient)
	}

	return patients, nil
}

// GetPatientInfo returns all patients info.
func (db DB) GetPatientsInfo(ctx context.Context) ([]dao.PatientInfo, error) {
	// get all patients info
	pkEx := expression.Key("pk").BeginsWith("PATIENT#")
	skEx := expression.Key("sk").Equal(expression.Value("INFO#0"))

	expr, err := expression.NewBuilder().
		WithKeyCondition(pkEx.And(skEx)).
		Build()
	if err != nil {
		return nil, err
	}

	res, err := db.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &db.Table,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		return []dao.PatientInfo{}, nil
	}

	patients := make([]dao.PatientInfo, 0, res.Count)
	for _, i := range res.Items {
		patient := dao.PatientInfo{}
		if err := attributevalue.UnmarshalMap(i, &patient); err != nil {
			return nil, err
		}

		count, err := db.countPatientDCM(ctx, patient.PK)
		if err != nil {
			return nil, err
		}

		patient.DCMCount = count

		patients = append(patients, patient)
	}

	return patients, nil
}

func (db DB) UpdatePatientInfo(ctx context.Context, pk string, data *dao.PatientInfo) error {
	data.UpdatedAt = time.Now()

	filterEx := expression.Name("pk").Equal(expression.Value(pk)).
		And(expression.Name("sk").Equal(expression.Value("INFO#0")))

	updateEx := expression.Set(expression.Name("filters"), expression.Value(data.Filters))
	updateEx.Set(expression.Name("lastname"), expression.Value(data.Lastname))
	updateEx.Set(expression.Name("firstname"), expression.Value(data.Firstname))
	updateEx.Set(expression.Name("updated_at"), expression.Value(data.UpdatedAt))

	updateExpr, err := expression.NewBuilder().
		WithUpdate(updateEx).
		Build()
	if err != nil {
		return err
	}

	filterExpr, err := expression.NewBuilder().
		WithFilter(filterEx).
		Build()
	if err != nil {
		return err
	}

	_, err = db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 &db.Table,
		ExpressionAttributeNames:  updateExpr.Names(),
		ExpressionAttributeValues: updateExpr.Values(),
		UpdateExpression:          updateExpr.Update(),
		Key:                       filterExpr.Values(),
	})

	return err
}

func (db DB) DeletePatient(ctx context.Context, pk string) error {
	filterEx := expression.Key("pk").Equal(expression.Value(pk))

	filterExpr, err := expression.NewBuilder().WithKeyCondition(filterEx).Build()
	if err != nil {
		return err
	}

	_, err = db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &db.Table,
		Key:       filterExpr.Values(),
	})

	return err
}

// countPatientDCM returns the number of processed DICOM files for a given patient.
func (db DB) countPatientDCM(ctx context.Context, pk string) (uint, error) {
	// count patient DCM
	pkEx := expression.Key("pk").Equal(expression.Value(pk))
	skEx := expression.Key("sk").BeginsWith("DCM#")

	expr, err := expression.NewBuilder().
		WithKeyCondition(pkEx.And(skEx)).
		Build()
	if err != nil {
		return 0, err
	}

	res, err := db.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &db.Table,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return 0, err
	}

	return uint(res.Count), nil
}
