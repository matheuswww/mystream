package upload_request

type Chunk struct {
	Hash string 		`json:"hash"`
	Chunk int 			`json:"chunk"`
	Data []byte 		`json:"data"`
}

type UploadFile struct {
	FileHash string `json:"fileHash"`
	Filename string `json:"fileName"`
	TotalChunk int 	`json:"totalChunk"`
	Chunks []Chunk 	`json:"chunks"`
}


