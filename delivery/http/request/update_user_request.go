package request

type UpdateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required,min=2,max=100"`
	LastName  string `json:"lastName" validate:"required,min=2,max=100"`
	IsAdmin   *bool  `json:"isAdmin" validate:"required"`
	Status    string `json:"status" validate:"required,oneof=active inactive suspended"`
	Password  string `json:"password" validate:"omitempty,password"`
}
