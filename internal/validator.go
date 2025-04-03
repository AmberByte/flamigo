package internal

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validate(obj interface{}) error {
	v := validator.New()
	v.RegisterValidation("is-username", func(fl validator.FieldLevel) bool {
		f := fl.Field()
		if f.Kind() != reflect.String {
			return false
		}
		valString := f.String()
		if len(valString) < 3 || len(valString) > 20 {
			return false
		}
		// This should alway compile so we can ignore the error
		usernameRegex, _ := regexp.Compile("^[a-zA-Z0-9_]+$")
		return usernameRegex.MatchString(valString)
	})
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
	err := v.Struct(obj)
	return err
}
