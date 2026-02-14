package dto_utils

type FileResponse struct {
	Id         string `json:"id,omitempty"`
	StorageKey string `json:"storage_key,omitempty"`
	MimeType   string `json:"mime_type,omitempty"`
	SizeBytes  int    `json:"size_bytes,omitempty"`
	Storage    string `json:"storage,omitempty"`
	FileUrl    string `json:"url,omitempty"`
}
