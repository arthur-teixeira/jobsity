package validators

type TaskRequest struct {
	Title string `json:"title"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
