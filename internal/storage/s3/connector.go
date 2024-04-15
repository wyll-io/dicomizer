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

// Upload sends "r" data to remote AWS "bucket" with key filename "key". "h" must be a base64 encoding sha256 hash
func (c *Client) Upload(ctx context.Context, r io.Reader, length int, h string, opts storage.Options) error {
	_, err := c.client.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:            &opts.Bucket,
		Key:               &opts.Key,
		Body:              r,
		// ChecksumSHA256:    &h,
		// ChecksumAlgorithm: "SHA256",
		ContentType:       aws.String("application/dicom"),
		ContentLength:     aws.Int64(int64(length)),
	})

	return err
}
