package models

const (
	FileStatusUploaded = "uploaded"
	FileStatusPending  = "pending"

	StorageS3    = "s3"
	StorageR2    = "r2"
	StorageLocal = "local"
)

type File struct {
	BaseModel
	AuthorID     int64
	StorageKey   string
	FileStatus   string
	OriginalName *string
	MimeType     *string
	SizeBytes    *int64
	Storage      string
}
