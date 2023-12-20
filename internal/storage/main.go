package storage

import (
	"context"
	"io"
)

type StorageAction interface {
	Upload(ctx context.Context, r io.Reader, opts interface{}) error
	UploadFile(ctx context.Context, fp string, opts interface{}) error
}
