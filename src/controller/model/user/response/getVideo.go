package user_response

type Url struct {
	FullHd string `json:"full_hd"`
	Hd     string `json:"hd"`
	Sd     string `json:"sd"`
	Vd     string `json:"vd"`
}

type GetVideo struct {
	Id 			    string `json:"id"`
	Title 			string `json:"title"`
	Description string `json:"description"`
	Url 				Url 	 `json:"url"`
	Cursor	 		string `json:"cursor"`
	FileHash    string `json:"-"`
}