package user

type RegisterUserInput struct {
	Name       string `json:"name" binding:"required"`
	Occupation string `json:"occupation" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,gte=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type EmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type FormCreateUserInput struct{
	Name	string `form:"name" bindung : "required"`
	Email	string `form:"email" bindung : "required,email"`
	Occupation	string `form:"occupation" bindung : "required"`
	Password	string `form:"password" bindung : "required"`
	Error error
}

type FormUpdateUserInput struct{
	ID		int
	Name	string `form:"name" bindung : "required"`
	Email	string `form:"email" bindung : "required,email"`
	Occupation	string `form:"occupation" bindung : "required"`
	Error error
}