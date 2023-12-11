package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/wyll-io/dicomizer/internal/models"
)

type DB struct {
	Client *dynamodb.Client
}

func New(cfg aws.Config) DB {
	return DB{
		Client: dynamodb.NewFromConfig(cfg),
	}
}

func (db DB) AddStudy(ctx context.Context, study models.Study) error {
	item, err := attributevalue.MarshalMap(study)
	if err != nil {
		return err
	}

	_, err = db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("studies"),
		Item:      item,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) AddPatient(ctx context.Context, patient models.Patient) error {
	item, err := attributevalue.MarshalMap(patient)
	if err != nil {
		return err
	}

	_, err = db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("patients"),
		Item:      item,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) GetPatient(ctx context.Context, uuid string, nestedValues bool) (models.Patient, []models.Study, error) {
	item, err := db.Client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"uuid": &types.AttributeValueMemberS{
				Value: uuid,
			},
		},
		TableName: aws.String("patients"),
	})
	if err != nil {
		return models.Patient{}, []models.Study{}, err
	}

	patient := models.Patient{}
	if err := attributevalue.UnmarshalMap(item.Item, &patient); err != nil {
		return models.Patient{}, []models.Study{}, err
	}

	if nestedValues {
		studies, err := db.GetStudiesByPatientUUID(ctx, patient.UUID)
		if err != nil {
			return models.Patient{}, []models.Study{}, err
		}

		return patient, studies, nil
	}

	return patient, []models.Study{}, nil
}

func (db DB) GetStudyByUUID(ctx context.Context, uuid string) (models.Study, error) {
	item, err := db.Client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"uuid": &types.AttributeValueMemberS{
				Value: uuid,
			},
		},
		TableName: aws.String("studies"),
	})
	if err != nil {
		return models.Study{}, err
	}

	study := models.Study{}
	if err := attributevalue.UnmarshalMap(item.Item, &study); err != nil {
		return models.Study{}, err
	}

	return study, nil
}

func (db DB) GetStudiesByPatientUUID(ctx context.Context, patientUUID string) ([]models.Study, error) {
	keyEx := expression.Key("patient_uuid").Equal(expression.Value(patientUUID))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return []models.Study{}, err
	}

	items, err := db.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String("studies"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return []models.Study{}, err
	}

	studies := make([]models.Study, 0, items.Count)
	attributevalue.UnmarshalListOfMaps(items.Items, &studies)

	return studies, nil
}
