package validation

import "testing"

var rules = Rules{
	"require":  Required,
	"length":   Length(15),
	"email":    Email,
	"password": Password(4),
}

func TestValidator(t *testing.T) {
	v := New(rules)

	t.Run("validates all fields and returns no errors", func(t *testing.T) {
		checks := map[string]Check{
			"name":     {Value: "John", RuleKey: []string{"require", "length"}},
			"email":    {Value: "test@email.com", RuleKey: []string{"require", "email", "length"}},
			"password": {Value: "Test1_", RuleKey: []string{"require", "password"}},
		}
		errs := v.Validate(checks)

		if errs != nil {
			t.Fatalf("unexpected errors returned validating data: %v", errs)
		}
	})

	t.Run("return 4 errors for invalid data validation", func(t *testing.T) {
		checks := map[string]Check{
			"name":     {Value: "", RuleKey: []string{"require"}},
			"drink":    {Value: "someGoodIceCoffee", RuleKey: []string{"length"}},
			"email":    {Value: "testemail.com", RuleKey: []string{"email"}},
			"password": {Value: "test", RuleKey: []string{"password"}},
		}

		errs := v.Validate(checks)

		if len(errs) != 4 {
			t.Fatalf("expects 4 validation errors got %d", len(errs))
		}
	})
}
