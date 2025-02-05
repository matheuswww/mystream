package upload_request

type UpdateVideo struct {
	FileHash    string `json:"file_hash" binding:"required"`
	Title 			string `json:"title"`
	Description string `json:"description"`
	Uploaded    *bool  `json:"uploaded"`
}