package upload_request

type FileHash struct {
	FileHash string `json:"file_hash" binding:"required"`
}