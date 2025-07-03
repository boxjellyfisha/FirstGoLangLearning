package hello

type GreetingResponse struct {
	Name    string  `json:"name"`
	Message *string `json:"message,omitempty"` // Optional string field
}
