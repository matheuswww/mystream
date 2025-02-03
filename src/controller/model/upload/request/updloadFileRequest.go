package upload_request

type Chunk struct {
	Hash string 		`json:"hash"`
	Chunk int 			`json:"chunk"`
	Data []byte 		`json:"data"`
}

type UploadFile struct {
	Title 	 		string 	`json:"title"`
	Description string 	`json:"description"`
	FileHash 		string 	`json:"fileHash"`
	Filename 		string 	`json:"fileName"`
	TotalChunk  int 		`json:"totalChunk"`
	Chunks 			[]Chunk `json:"chunks" validate:"dive"`
}


