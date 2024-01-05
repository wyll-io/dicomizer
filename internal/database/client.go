package database

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dao "github.com/wyll-io/dicomizer/internal/DAO"
)

type DB struct {
	Client *dynamodb.Client
}

type UUIDFilter struct {
	UUID string
	Key  string
}

func New(cfg aws.Config) dao.DBActions {
	return DB{
		Client: dynamodb.NewFromConfig(cfg),
	}
}

func (db DB) add(ctx context.Context, collection string, data interface{}) error {
	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		return err
	}

	_, err = db.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(collection),
		Item:      item,
	})
	if err != nil {
		return err
	}

	return nil
}

func (db DB) update(ctx context.Context, collection, uuid string, data interface{}) error {
	item, err := attributevalue.MarshalMap(data)
	if err != nil {
		return err
	}

	update := expression.UpdateBuilder{}
	for k, v := range item {
		update = update.Set(expression.Name(k), expression.Value(v))
	}

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil
	}

	_, err = db.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"uuid": &types.AttributeValueMemberS{
				Value: uuid,
			},
		},
		TableName:                 aws.String(collection),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	})

	return err
}

func (db DB) getByUUID(ctx context.Context, collection, uuid string, out interface{}) error {
	item, err := db.Client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"uuid": &types.AttributeValueMemberS{
				Value: uuid,
			},
		},
		TableName: aws.String(collection),
	})
	if err != nil {
		return err
	}

	return attributevalue.UnmarshalMap(item.Item, &out)
}

func (db DB) getByParentUUID(ctx context.Context, collection string, filter UUIDFilter, out interface{}) error {
	keyEx := expression.Key(filter.UUID).Equal(expression.Value(filter.Key))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return err
	}

	items, err := db.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(collection),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return err
	}

	return attributevalue.UnmarshalListOfMaps(items.Items, &out)
}

func (db DB) AddStudy(ctx context.Context, study dao.DCMImage) error {
	return db.add(ctx, "studies", ConvertDAOToStudy(study))
}

func (db DB) AddPatient(ctx context.Context, patient dao.Patient) error {
	p, sts := ConvertDAOToPatient(patient)
	if err := db.add(ctx, "patients", p); err != nil {
		return err
	}

	if len(sts) > 0 {
		for _, s := range sts {
			if err := db.add(ctx, "studies", s); err != nil {
				return err
			}
		}
	}

	return nil
}

func (db DB) GetPatientByUUID(ctx context.Context, uuid string, nestedValues bool) (dao.Patient, error) {
	p := patient{}

	if err := db.getByUUID(ctx, "patients", uuid, &p); err != nil {
		return dao.Patient{}, err
	}

	sts := []dcmImage{}
	if nestedValues {
		err := db.getByParentUUID(ctx, "studies", UUIDFilter{UUID: uuid, Key: "patient_uuid"}, &sts)
		if err != nil {
			return dao.Patient{}, err
		}
	}

	return ConvertPatientToDAO(p, sts), nil
}
func (db DB) GetPatient(ctx context.Context, filters dao.SearchPatientParams, nestedValues bool) ([]dao.Patient, error) {
	expr := expression.Contains(expression.Name("firstname"), filters.Firstname).
		Or(expression.Contains(expression.Name("lastname"), filters.Lastname)).
		Or(expression.Contains(expression.Name("filters"), filters.Filters))
	r, err := expression.NewBuilder().WithCondition(expr).Build()
	if err != nil {
		return []dao.Patient{}, err
	}

	patients := []patient{}

	out, err := db.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String("patients"),
		KeyConditionExpression:    r.KeyCondition(),
		ExpressionAttributeNames:  r.Names(),
		ExpressionAttributeValues: r.Values(),
	})
	if err != nil {
		return []dao.Patient{}, err
	}

	if err := attributevalue.UnmarshalListOfMaps(out.Items, &patients); err != nil {
		return []dao.Patient{}, err
	}

	studies := map[string][]dcmImage{}
	if nestedValues {
		for _, p := range patients {
			sts := []dcmImage{}
			err := db.getByParentUUID(ctx, "studies", UUIDFilter{UUID: p.UUID, Key: "patient_uuid"}, &sts)
			if err != nil {
				return []dao.Patient{}, err
			}
			studies[p.UUID] = sts
		}
	}

	return ConvertPatientsToDAO(patients, studies), nil
}

func (db DB) GetStudyByUUID(ctx context.Context, uuid string) (dao.DCMImage, error) {
	s := dcmImage{}

	if err := db.getByUUID(ctx, "studies", uuid, &s); err != nil {
		return dao.DCMImage{}, err
	}

	return ConvertStudyToDAO(s), nil
}

func (db DB) GetStudiesByPatientUUID(ctx context.Context, patientUUID string) ([]dao.DCMImage, error) {
	sts := []dcmImage{}

	err := db.getByParentUUID(ctx, "studies", UUIDFilter{UUID: patientUUID, Key: "patient_uuid"}, &sts)
	if err != nil {
		return nil, err
	}

	return ConvertStudiesToDAO(sts), nil
}

func (db DB) DeletePatient(ctx context.Context, uuid string) error {
	_, err := db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String("patients"),
		Key: map[string]types.AttributeValue{
			"uuid": &types.AttributeValueMemberS{
				Value: uuid,
			},
		},
	})
	if err != nil {
		return err
	}

	// TODO: delete all studies related to this patient

	return nil
}

func (db DB) UpdatePatient(ctx context.Context, patient dao.Patient) error {
	patient.UpdatedAt = time.Now()
	p, sts := ConvertDAOToPatient(patient)
	if err := db.update(ctx, "patients", p.UUID, p); err != nil {
		return err
	}

	if len(sts) > 0 {
		for _, s := range sts {
			s.UpdatedAt = time.Now()
			if err := db.update(ctx, "studies", s.UUID, s); err != nil {
				return err
			}
		}
	}

	return nil
}
