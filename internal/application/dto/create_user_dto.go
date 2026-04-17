package dto

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required,min=2,max=100"`
	LastName  string `json:"lastName" validate:"required,min=2,max=100"`
	IsAdmin   *bool  `json:"isAdmin" validate:"required"`
	Password  string `json:"password" validate:"required,password"`
}
