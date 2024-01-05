package database

import dao "github.com/wyll-io/dicomizer/internal/DAO"

func ConvertPatientToDAO(p patient, sts []dcmImage) dao.Patient {
	patientDAO := dao.Patient{
		UUID:      p.UUID,
		Firstname: p.Firstname,
		Lastname:  p.Lastname,
		Filters:   p.Filters,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,
		Studies:   make([]dao.DCMImage, 0, len(sts)),
	}

	for _, s := range sts {
		patientDAO.Studies = append(patientDAO.Studies, ConvertStudyToDAO(s))
	}

	return patientDAO
}

func ConvertStudyToDAO(s dcmImage) dao.DCMImage {
	return dao.DCMImage{
		UUID:        s.UUID,
		PatientUUID: s.PatientUUID,
		Hash:        s.Hash,
		Filename:    s.Filename,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		DeletedAt:   s.DeletedAt,
	}
}

func ConvertStudiesToDAO(sts []dcmImage) []dao.DCMImage {
	studiesDAO := make([]dao.DCMImage, 0, len(sts))

	for _, s := range sts {
		studiesDAO = append(studiesDAO, ConvertStudyToDAO(s))
	}

	return studiesDAO
}

func ConvertPatientsToDAO(ps []patient, sts map[string][]dcmImage) []dao.Patient {
	psDAO := make([]dao.Patient, 0, len(ps))
	for _, p := range ps {
		sts := sts[p.UUID]
		psDAO = append(psDAO, ConvertPatientToDAO(p, sts))
	}

	return psDAO
}

func ConvertDAOToPatient(p dao.Patient) (patient, []dcmImage) {
	return patient{
		UUID:      p.UUID,
		Firstname: p.Firstname,
		Lastname:  p.Lastname,
		Filters:   p.Filters,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,
	}, ConvertDAOToStudies(p.Studies)
}

func ConvertDAOToStudies(studies []dao.DCMImage) []dcmImage {
	studiesModels := make([]dcmImage, 0, len(studies))

	for _, s := range studies {
		studiesModels = append(studiesModels, ConvertDAOToStudy(s))
	}

	return studiesModels
}

func ConvertDAOToStudy(s dao.DCMImage) dcmImage {
	return dcmImage{
		UUID:        s.UUID,
		PatientUUID: s.PatientUUID,
		Hash:        s.Hash,
		Filename:    s.Filename,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		DeletedAt:   s.DeletedAt,
	}
}
