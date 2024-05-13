package s3

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/wyll-io/dicomizer/internal/storage"
)

type Options struct {
	Bucket string
	Key    string
}

func (c *Client) UploadFile(ctx context.Context, fp string, opts storage.Options) error {
	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	body := bytes.NewBuffer(nil)
	d := io.TeeReader(f, body)

	h := sha256.New()
	wr, err := io.Copy(h, d)
	if err != nil {
		return err
	}

	_, err = c.client.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:         &opts.Bucket,
		Key:            &opts.Key,
		Body:           body,
		ChecksumSHA256: aws.String(base64.StdEncoding.EncodeToString(h.Sum(nil))),
		ContentLength:  aws.Int64(wr),
		ContentType:    aws.String("application/dicom"),
	})

	return err
}

// Upload sends "r" data to remote AWS "bucket" with key filename "key"
func (c *Client) Upload(ctx context.Context, r io.Reader, length int, opts storage.Options) error {
	body := bytes.NewBuffer(nil)
	d := io.TeeReader(r, body)

	h := sha256.New()
	wr, err := io.Copy(h, d)
	if err != nil {
		return err
	}

	_, err = c.client.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:         &opts.Bucket,
		Key:            &opts.Key,
		Body:           body,
		ChecksumSHA256: aws.String(base64.StdEncoding.EncodeToString(h.Sum(nil))),
		ContentType:    aws.String("application/dicom"),
		ContentLength:  aws.Int64(wr),
	})

	return err
}

func (c *Client) DeletePatientFiles(ctx context.Context, pk_key string) error {
	fmt.Printf("Deleting patient files in %s...\n", pk_key)
	rsp, err := c.client.ListObjectsV2(ctx, &awsS3.ListObjectsV2Input{
		Bucket: &c.bucket,
		Prefix: &pk_key,
	})
	if err != nil {
		return err
	}

	for _, obj := range rsp.Contents {
		_, err := c.client.DeleteObject(ctx, &awsS3.DeleteObjectInput{
			Bucket: &c.bucket,
			Key:    obj.Key,
		})
		if err != nil {
			return err
		}

	}

	fmt.Println("Patient's files deleted successfully")

	return nil
}
