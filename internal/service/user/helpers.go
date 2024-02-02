package user

import (
	"fmt"
	"inditilla/internal/entity"
	"inditilla/internal/service/validator"
	"regexp"
)

const (
	maxInitialsLen = 255
	maxEmailLen    = 255
	minPasswordLen = 8
	maxPasswordLen = 500
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9]{0, 61}[a-zA-Z0-9])?)*$")

func isRightSignUp(u *entity.UserSignupForm) bool {
	u.CheckField(validator.NotBlank(u.FirstName), "firstName", "This field cannot be blank")
	u.CheckField(validator.MaxChar(u.FirstName, maxInitialsLen), "firstName", fmt.Sprintf("Maximum characters length exceeded - %d", maxInitialsLen))
	u.CheckField(validator.NotBlank(u.LastName), "lastName", "This field cannot be blank")
	u.CheckField(validator.MaxChar(u.LastName, maxInitialsLen), "lastName", fmt.Sprintf("Maximum characters length exceeded - %d", maxInitialsLen))
	u.CheckField(validator.NotBlank(u.Email), "email", "This field cannot be blank")
	u.CheckField(validator.MaxChar(u.Email, maxEmailLen), "email", fmt.Sprintf("Maximum characters length exceeded - %d", maxEmailLen))
	u.CheckField(validator.Matches(u.Email, EmailRX), "email", "Invalid email address")
	u.CheckField(validator.NotBlank(u.Password), "password", "This field cannot be blank")
	u.CheckField(validator.MinChar(u.Password, minPasswordLen), "password", fmt.Sprintf("This field should be %d characters length minimum", minPasswordLen))
	u.CheckField(validator.MaxChar(u.Password, maxPasswordLen), "password", fmt.Sprintf("Maximum characters length exceeded - %d", maxPasswordLen))

	return u.Valid()
}

func isRightLogin(u *entity.UserLoginForm) bool {
	u.CheckField(validator.NotBlank(u.Email), "email", "This field cannot be blank")
	u.CheckField(validator.Matches(u.Email, EmailRX), "email", "This should valid email address")
	u.CheckField(validator.NotBlank(u.Password), "password", "This field cannot be blank")

	return u.Valid()
}
