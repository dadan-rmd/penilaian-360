package requestvalidationerror

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/stoewer/go-strcase"
)

type ValidationField struct {
	Field   string
	Message string
}

var (
	RequiredMsg = "is required"
	MaxMsg      = "maximum : "
	MinMsg      = "minimum : "
	NoChange    = "no change"
)

func GetvalidationError(err error) []ValidationField {
	var validationFields []ValidationField
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, validationError := range ve {

			switch validationError.Tag() {
			case "required":
				myField := validationError.Field()
				validationFields = append(validationFields, ValidationField{
					Field:   strcase.LowerCamelCase(myField),
					Message: "this " + RequiredMsg,
				})
			case "max":
				myField := validationError.Field()
				validationFields = append(validationFields, ValidationField{
					Field:   strcase.LowerCamelCase(myField),
					Message: MaxMsg + validationError.Param(),
				})
			case "min":
				myField := validationError.Field()
				validationFields = append(validationFields, ValidationField{
					Field:   strcase.LowerCamelCase(myField),
					Message: MinMsg + validationError.Param(),
				})
			}
		}
	}

	return validationFields
}

func GetvalidationError2(err error) error {
	arr := GetvalidationError(err)
	arrMsg := []string{}
	for _, v := range arr {
		arrMsg = append(arrMsg, fmt.Sprintf("%s: %s", v.Field, v.Message))
	}

	if len(arrMsg) > 0 {
		return fmt.Errorf(strings.Join(arrMsg, ", "))
	}

	return nil
}

func IsINAPhoneValid(phone string) bool {
	rgx, _ := regexp.Compile(`^(\+62)[0-9]{10,12}$`)

	return rgx.MatchString(phone)
}

func IsEmailValid(email string) bool {
	rgx, _ := regexp.Compile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

	return rgx.MatchString(email)
}
