package admin_request

type Signin struct {
	Email 	 string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}