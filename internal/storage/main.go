package storage

import (
	"context"
	"io"
)

type Options struct {
	// S3: Bucket is the name of the bucket to upload to
	Bucket string
	// S3: Key is the name of the file to upload
	Key string

	// Glacier: AccountID is the AWS account ID of the account that owns the vault
	AccountID string
	// Glacier: ArchiveDescription is a description of the archive
	ArchiveDescription string
	// Glacier: VaultName is the name of the vault
	VaultName string
}

type StorageAction interface {
	Upload(ctx context.Context, r io.Reader, length int, h string, opts Options) error
	UploadFile(ctx context.Context, fp string, opts Options) error
}
