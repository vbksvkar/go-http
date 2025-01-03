package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func NewValidator(tagNameToRegister string) *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get(tagNameToRegister), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}
