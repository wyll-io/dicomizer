package s3

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
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

func (c *Client) Upload(ctx context.Context, r io.Reader, opts storage.Options) error {
	h := sha256.New()
	wr, err := io.Copy(h, r)
	if err != nil {
		return err
	}

	c.client.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:            &opts.Bucket,
		Key:               &opts.Key,
		Body:              r,
		ChecksumSHA256:    aws.String(string(h.Sum(nil))),
		ChecksumAlgorithm: "SHA256",
		ContentLength:     aws.Int64(wr),
		ContentType:       aws.String("application/dicom"),
	})

	return nil
}
