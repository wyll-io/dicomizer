package check

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
	"path/filepath"

	"github.com/suyashkumar/dicom"
	dao "github.com/wyll-io/dicomizer/internal/DAO"
	"github.com/wyll-io/dicomizer/internal/storage"
	"github.com/wyll-io/dicomizer/pkg/anonymize"
	"github.com/wyll-io/dicomizer/pkg/network"
)

func CheckPatientDCM(
	ctx context.Context,
	storageClient storage.StorageAction,
	dbClient dao.DBActions,
	pacs, aet, aec, aem string,
	pInfo dao.PatientInfo,
) error {
	fmt.Println("Fetching patients DCM files...")
	tmp, err := network.MoveSCU(pacs, aet, aec, aem, pInfo.Filters)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(tmp)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Processing %d DCM files...\n", len(files))
	for _, f := range files {
		if f.IsDir() {
			panic("unexpected directory")
		}

		fmt.Println(f.Name())
		if strings.Contains(f.Name(), "rsp") {
		  fmt.Println("Skipping query file...")
		  continue
		}

		dataset, err := anonymizeDataset(filepath.Join(tmp, f.Name()))
		if err != nil {
			fmt.Printf("failed to anonymize: %s", f.Name())
			return err
		}

		outF, err := os.OpenFile(filepath.Join(tmp, fmt.Sprintf("%s.anonymized", f.Name())), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		defer outF.Close()

		if err := dicom.Write(outF, dataset); err != nil {
			fmt.Printf("failed to create anonymized file: %s", f.Name())
			return err
		}

		h, err := getHash(outF)
		if err != nil {
			return err
		}

		found, err := dbClient.CheckDCM(ctx, h, f.Name())
		if err != nil {
			return err
		}
		if found {
			fmt.Printf("DCM file found in DB, skipping: %s\n", f.Name())
			continue
		}
		fmt.Printf("DCM file not found in DB, uploading: %s\n", f.Name())
		err = processFoundDCM(ctx, storageClient, dbClient, filepath.Join(tmp, f.Name()), h, pInfo.PK)
		if err != nil {
			return err
		}
	}

	fmt.Println("Patient processed. Cleaning up...")

	return os.RemoveAll(tmp)
}

func getHash(r io.Reader) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, r); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func processFoundDCM(
	ctx context.Context,
	storageClient storage.StorageAction,
	dbClient dao.DBActions,
	fp, hash, pk string,
) error {
	fname := filepath.Base(fp)
	if err := storageClient.UploadFile(ctx, fp, storage.Options{
		Bucket: "dicomizer",
		Key:    fname,
	}); err != nil {
		return err
	}

	return dbClient.AddDCM(ctx, pk, &dao.DCMInfo{
		Hash:     hash,
		Filename: fname,
	})
}

func anonymizeDataset(fp string) (dicom.Dataset, error) {
	dataset, err := dicom.ParseFile(fp, nil)
	if err != nil {
		return dicom.Dataset{}, err
	}

	if err := anonymize.AnonymizeDataset(&dataset); err != nil {
		return dicom.Dataset{}, err
	}

	return dataset, nil
}
