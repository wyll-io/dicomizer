package glacier

import (
	"context"
	"crypto/sha256"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsG "github.com/aws/aws-sdk-go-v2/service/glacier"
)

type Options struct {
	AccountID          string
	ArchiveDescription string
	VaultName          string
}

func (c *Client) UploadFile(ctx context.Context, fp string, params interface{}) error {
	opts := params.(Options)

	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	input := &awsG.UploadArchiveInput{
		AccountId:          &opts.AccountID,
		ArchiveDescription: &opts.ArchiveDescription,
		Body:               f,
		Checksum:           aws.String(string(h.Sum(nil))),
		VaultName:          &opts.VaultName,
	}

	_, err = c.client.UploadArchive(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Upload(ctx context.Context, r io.Reader, params interface{}) error {
	opts := params.(Options)

	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return err
	}

	input := &awsG.UploadArchiveInput{
		AccountId:          &opts.AccountID,
		ArchiveDescription: &opts.ArchiveDescription,
		Body:               r,
		Checksum:           aws.String(string(h.Sum(nil))),
		VaultName:          &opts.VaultName,
	}

	_, err := c.client.UploadArchive(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
