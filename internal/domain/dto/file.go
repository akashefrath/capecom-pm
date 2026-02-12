package dto

type CreateFileRequest struct {
	FileName string `json:"file_name" form:"file_name" binding:"required"`
	MimeType string `json:"mime_type" form:"mime_type" binding:"required"`
}

type CreateFileResponse struct {
	FileID    string `json:"file_id"`
	UploadURL string `json:"upload_url"`
}

type ConfirmUploadRequest struct {
	FileIDs []string `json:"file_ids" binding:"required,min=1"`
}

type FileUploadStatus struct {
	FileID  string `json:"file_id"`
	Exists  bool   `json:"exists"`
	Updated bool   `json:"updated"`
}
