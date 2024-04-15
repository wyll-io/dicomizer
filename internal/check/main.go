package check

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	pacs, aet, aec, aem, labo string,
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
			panic("unexpected directory")
		}

		if strings.Contains(f.Name(), "rsp") {
			fmt.Println("Skipping query file...")
			continue
		}

		dataset, err := anonymizeDataset(filepath.Join(tmp, f.Name()))
		if err != nil {
			return fmt.Errorf("failed to anonymize dataset: %v", err)
		}

		data := bytes.NewBuffer(nil)
		if err := dicom.Write(data, dataset); err != nil {
			return fmt.Errorf("failed to write anonymized dataset: %v", err)
		}

		err = func() error {
			h, err := getHash(data)
			if err != nil {
				return err
			}

      fileKey := fmt.Sprintf("%s/%s/%s", labo, strings.Replace(pInfo.PK, "PATIENT#", "", 1), f.Name())
      found, err := dbClient.CheckDCM(ctx, hex.EncodeToString(h), fileKey)
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
				h,
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

func getHash(r io.Reader) ([]byte, error) {
	// hasher := sha256.New()
	// if _, err := io.Copy(hasher, r); err != nil {
	// 	return nil, fmt.Errorf("failed to calculate sha256 hash for anonymized dataset: %v", err)
	// }
	//
	// return hasher.Sum(nil), nil

  return []byte("disabled for debugging"), nil
}

func processFoundDCM(
	ctx context.Context,
	storageClient storage.StorageAction,
	dbClient dao.DBActions,
	fileKey, pk string,
	h []byte,
	b *bytes.Buffer,
) error {
	if err := storageClient.Upload(ctx, b, b.Len(), string(h), storage.Options{
		Bucket: "dicomizer",
		Key:    fileKey,
	}); err != nil {
		return err
	}

	return dbClient.AddDCM(ctx, pk, &dao.DCMInfo{
		Hash:     hex.EncodeToString(h),
		Filename: fileKey,
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
