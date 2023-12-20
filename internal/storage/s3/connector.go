package s3

import (
	"context"
	"crypto/sha256"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type Options struct {
	Bucket string
	Key    string
}

func (c *Client) UploadFile(ctx context.Context, fp string, params interface{}) error {
	opts := params.(Options)

	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	wr, err := io.Copy(h, f)
	if err != nil {
		return err
	}

	c.client.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:            &opts.Bucket,
		Key:               &opts.Key,
		Body:              f,
		ChecksumSHA256:    aws.String(string(h.Sum(nil))),
		ChecksumAlgorithm: "SHA256",
		ContentLength:     aws.Int64(wr),
		ContentType:       aws.String("application/dicom"),
	})

	return nil
}

func (c *Client) Upload(ctx context.Context, r io.Reader, params interface{}) error {
	opts := params.(Options)

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
