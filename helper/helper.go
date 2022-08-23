/*
untuk format response json
yg terdiri dari object meta dan data
*/
package helper

import (
	// "fmt"
	// "fmt"
	"math/rand"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"` //berbentuk meta
	Data interface{} `json:"data"` //berbentuk interface, karena fleksibel
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}


func FormatValidationError(err error) []string {
	var errors []string
	// fmt.Print(err.Error());
	for _, e := range err.(validator.ValidationErrors) {
		// errors = append(errors, e.Error())
		// fmt.Println(e.Field())
		if e.Tag() == "required"{
			errors = append(errors, e.Field() + " harus diisi")
		}
		if e.Tag() == "email"{
			errors = append(errors, "Format "+e.Field()+" tidak valid")
		}
		if e.Field() == "Password" && e.Tag() == "gte"{
			errors = append(errors, e.Field() + " harus lebih dari 6 karakter")
		}
		// fmt.Println(e.Tag())
	}

	return errors
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ#@$%^|~*")

func TokenString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
