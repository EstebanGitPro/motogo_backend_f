package output

import "context"

// StorageFileDeleter is the interface for deleting files from cloud storage (Firebase Storage).
// Named per Go convention: single-method interfaces use the method verb + "er" suffix.
type StorageFileDeleter interface {
	// DeleteStorageFile deletes a file from cloud storage given its URL
	DeleteStorageFile(ctx context.Context, fileURL string) error
}
