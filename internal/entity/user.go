package entity

type UserSignupForm struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UserLoginForm struct {
	Identifier string
	Password   string
	validator.Validator
}
