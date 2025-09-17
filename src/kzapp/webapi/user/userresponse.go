package user

type UserResponse struct {
	Name    string  `json:"name"`
	Message *string `json:"message,omitempty"` // Optional string field
}


type UserUpdateRequest struct {
	Name    string  `json:"name" default:"your name"`
	Email *string `json:"email,omitempty" default:"name@mail.com"` // Optional string field
	Password *string `json:"password,omitempty" default:"***"` // Optional string field
}


// {
// 	"_id": {
// 	  "$oid": "687df826971b74d73ece63b6"
// 	},
// 	"name": "test",
// 	"email": "test@test.com",
// 	"password": "test"
//   }