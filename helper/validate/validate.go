package validate

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
	"github.com/mohdjishin/SplitWise/config"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("password_complexity", passwordComplexity)
	_ = validate.RegisterValidation("dateFormat", validateDateFormat)
}

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			tag := err.Tag()

			switch tag {
			case "required":
				return fmt.Errorf("field '%s' is required", fieldName)

			case "email":
				return fmt.Errorf("field '%s' must be a valid email address", fieldName)
			case "password_complexity":
				return fmt.Errorf("field '%s' must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one digit and one special character", fieldName)
			default:
				return fmt.Errorf("field '%s' failed validation for an unknown reason", fieldName)
			}
		}
	}

	return nil
}

func passwordComplexity(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	if config.GetConfig().ENV != "dev" { // Skip complexity check in dev mode
		if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
			return false
		}
		if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
			return false
		}
		if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
			return false
		}

		if matched, _ := regexp.MatchString(`[\W_]`, password); !matched {
			return false
		}
	}
	return true
}

func validateDateFormat(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() { // Check if the field is empty
		return true // If empty, return true (valid)
	}
	dateStr := fl.Field().String()
	// Regular expression to match YYYY-MM-DD format
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return re.MatchString(dateStr)
}
