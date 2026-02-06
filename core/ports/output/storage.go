package output

import "context"

// StorageClient interface for file storage operations (Firebase Storage)
type StorageClient interface {
	// DeleteStorageFile deletes a file from cloud storage given its URL
	DeleteStorageFile(ctx context.Context, fileURL string) error
}
