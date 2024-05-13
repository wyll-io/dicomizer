package check

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
	dao "github.com/wyll-io/dicomizer/internal/DAO"
	"github.com/wyll-io/dicomizer/internal/storage"
	"github.com/wyll-io/dicomizer/pkg/anonymize"
	"github.com/wyll-io/dicomizer/pkg/network"
)

func CheckPatientDCM(
	ctx context.Context,
	storageClient storage.StorageAction,
	dbClient dao.DBActions,
	pacs, aet, aec, aem, center string,
	pInfo dao.PatientInfo,
) error {
	fmt.Println("Fetching patients DCM files...")
	tmp, err := network.MoveSCU(pacs, aet, aec, aem, pInfo.Filters)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(tmp)
	if err != nil {
		return err
	}

	fmt.Printf("Processing %d DCM files...\n", len(files))
	for _, f := range files {
		if f.IsDir() {
			fmt.Println("unexpected directory")
			os.Exit(1)
		}

		if strings.Contains(f.Name(), "rsp") {
			fmt.Println("Skipping query file...")
			continue
		}

		dataset, examDate, examType, err := anonymizeDataset(filepath.Join(tmp, f.Name()))
		if err != nil {
			return fmt.Errorf("failed to anonymize dataset: %v", err)
		}

		data := bytes.NewBuffer(nil)
		if err := dicom.Write(data, dataset); err != nil {
			return fmt.Errorf("failed to write anonymized dataset: %v", err)
		}

		err = func() error {
			fileKey := fmt.Sprintf(
				"%s/%s/%s-%s/%s",
				center,
				strings.Replace(pInfo.PK, "PATIENT#", "", 1),
				examDate,
				examType,
				f.Name(),
			)
			found, err := dbClient.CheckDCM(ctx, fileKey)
			if err != nil {
				return fmt.Errorf("failed to check if DCM exists in DB: %v", err)
			}
			if found {
				fmt.Printf("DCM file found in DB, skipping: %s\n", f.Name())
				return nil
			}

			fmt.Printf("DCM file not found in DB, uploading: %s\n", f.Name())
			if err := processFoundDCM(
				ctx,
				storageClient,
				dbClient,
				fileKey,
				pInfo.PK,
				data,
			); err != nil {
				return fmt.Errorf("failed to upload anonymized DCM file: %v", err)
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}

	fmt.Println("Patient processed. Cleaning up...")
	return os.RemoveAll(tmp)
}

func processFoundDCM(
	ctx context.Context,
	storageClient storage.StorageAction,
	dbClient dao.DBActions,
	fileKey, pk string,
	b *bytes.Buffer,
) error {
	if err := storageClient.Upload(ctx, b, b.Len(), storage.Options{
		Bucket: "dicomizer",
		Key:    fileKey,
	}); err != nil {
		return err
	}

	return dbClient.AddDCM(ctx, pk, &dao.DCMInfo{
		Filename: fileKey,
	})
}

// anonymizeDataset read input file, convert it to a DICOM dataset and anonymize it.
// It returns the anonymized dataset, its study date and its examen type
func anonymizeDataset(fp string) (dicom.Dataset, string, string, error) {
	dataset, err := dicom.ParseFile(fp, nil)
	if err != nil {
		return dicom.Dataset{}, "", "", err
	}

	etTag, err := dataset.FindElementByTag(tag.Tag{
		Group:   0x0008,
		Element: 0x0060,
	})
	if err != nil {
		return dicom.Dataset{}, "", "", err
	}
	examType := etTag.Value.String()

	sdTag, err := dataset.FindElementByTag(tag.Tag{
		Group:   0x0008,
		Element: 0x0020,
	})
	if err != nil {
		return dicom.Dataset{}, "", "", err
	}
	studyDate := sdTag.Value.String()

	if err := anonymize.AnonymizeDataset(&dataset); err != nil {
		return dicom.Dataset{}, "", "", err
	}

	return dataset, studyDate, examType, nil
}
