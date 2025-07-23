package user

type UserResponse struct {
	Name    string  `json:"name"`
	Message *string `json:"message,omitempty"` // Optional string field
}
