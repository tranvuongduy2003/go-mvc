package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

func InitValidator() error {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	english := en.New()
	uni := ut.New(english, english)
	var found bool
	translator, found = uni.GetTranslator("en")
	if !found {
		return fmt.Errorf("translator not found")
	}

	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		return fmt.Errorf("failed to register translations: %w", err)
	}

	registerCustomValidations()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		*v = *validate
	}

	return nil
}

func ValidateStruct(data interface{}) map[string]string {
	errors := make(map[string]string)

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Translate(translator)
		}
	}

	return errors
}

func registerCustomValidations() {
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return len(password) >= 8 &&
			strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") &&
			strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") &&
			strings.ContainsAny(password, "0123456789") &&
			strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?")
	})

	validate.RegisterTranslation("password", translator, func(ut ut.Translator) error {
		return ut.Add("password", "{0} must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})

	validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		return len(phone) >= 10 && len(phone) <= 15 && strings.HasPrefix(phone, "+")
	})

	validate.RegisterTranslation("phone", translator, func(ut ut.Translator) error {
		return ut.Add("phone", "{0} must be a valid phone number starting with + and 10-15 digits", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})
}
