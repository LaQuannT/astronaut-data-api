package validation

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
)

type (
	Rule func(key string, value interface{}) error

	Rules map[string]Rule

	Validator struct {
		rules Rules
	}

	Check struct {
		Value   interface{}
		RuleKey []string
	}
)

var emailRegex = regexp.MustCompile(`[\w-\.]+@([\w-]+\.)+[\w-]{2,4}`)

func New(r Rules) *Validator {
	return &Validator{rules: r}
}

func (v *Validator) Validate(data map[string]Check) []error {
	var errs []error

	for key, check := range data {
		for _, rk := range check.RuleKey {
			rule := v.rules[rk]
			if err := rule(key, check.Value); err != nil {
				errs = append(errs, err)
			}

		}
	}

	return errs
}

func Required(key string, value interface{}) error {
	if value == "" {
		return fmt.Errorf("%s must not be blank", key)
	}
	return nil
}

func Length(max int) Rule {
	return func(key string, value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("%s is not a string", key)
		}
		if len(str) > max {
			return fmt.Errorf("%s must be less than or equal to %d characters", key, max)
		}
		return nil
	}
}

func Email(key string, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s is not a string", key)
	}

	if !emailRegex.MatchString(str) {
		return fmt.Errorf("%s is not a vaild email address", key)
	}
	return nil
}

func Password(min int) Rule {
	return func(key string, value interface{}) error {
		// password must contain atleast one uppercase, one lowercase, one number, one special character
		// and min length
		var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool

		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("%s is not a string", key)
		}

		if len(str) >= min {
			hasMinLen = true
		}

		for _, char := range str {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsNumber(char):
				hasNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}

		if !hasMinLen || !hasUpper || !hasLower || !hasNumber || !hasSpecial {
			return fmt.Errorf("%s must contain one upper, lower, number, and special character atleast %d characters", key, min)
		}
		return nil
	}
}

func Role(key string, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s is not a string", key)
	}

	if str == model.AdminUser || str == model.BaseUser {
		return nil
	} else {
		return fmt.Errorf("%s is not a valid role", key)
	}
}
