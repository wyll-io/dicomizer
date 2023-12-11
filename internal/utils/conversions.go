package utils

import (
	dao "github.com/wyll-io/dicomizer/internal/DAO"
	"github.com/wyll-io/dicomizer/internal/models"
)

func ConvertPatientToDAO(patient models.Patient, studies []models.Study) dao.Patient {
	patientDAO := dao.Patient{
		UUID:      patient.UUID,
		Firstname: patient.Firstname,
		Lastname:  patient.Lastname,
		Filters:   patient.Filters,
		CreatedAt: patient.CreatedAt,
		UpdatedAt: patient.UpdatedAt,
		DeletedAt: patient.DeletedAt,
		Studies:   make([]dao.Study, 0, len(studies)),
	}

	for _, s := range studies {
		patientDAO.Studies = append(patientDAO.Studies, ConvertStudyToDAO(s))
	}

	return patientDAO
}

func ConvertStudyToDAO(study models.Study) dao.Study {
	return dao.Study{
		UUID:        study.UUID,
		PatientUUID: study.PatientUUID,
		Status:      study.Status,
		Hash:        study.Hash,
		Filename:    study.Filename,
		CreatedAt:   study.CreatedAt,
		UpdatedAt:   study.UpdatedAt,
		DeletedAt:   study.DeletedAt,
	}
}

func ConvertDAOToPatient(patient dao.Patient) (models.Patient, []models.Study) {
	return models.Patient{
		UUID:      patient.UUID,
		Firstname: patient.Firstname,
		Lastname:  patient.Lastname,
		Filters:   patient.Filters,
		CreatedAt: patient.CreatedAt,
		UpdatedAt: patient.UpdatedAt,
		DeletedAt: patient.DeletedAt,
	}, ConvertDAOToStudies(patient.Studies)
}

func ConvertDAOToStudies(studies []dao.Study) []models.Study {
	studiesModels := make([]models.Study, 0, len(studies))

	for _, s := range studies {
		studiesModels = append(studiesModels, ConvertDAOToStudy(s))
	}

	return studiesModels
}

func ConvertDAOToStudy(study dao.Study) models.Study {
	return models.Study{
		UUID:        study.UUID,
		PatientUUID: study.PatientUUID,
		Status:      study.Status,
		Hash:        study.Hash,
		Filename:    study.Filename,
		CreatedAt:   study.CreatedAt,
		UpdatedAt:   study.UpdatedAt,
		DeletedAt:   study.DeletedAt,
	}
}
