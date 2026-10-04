package validation

import "github.com/go-playground/validator/v10"

// Register adds struct-level validation for T that reports only Fatal violations,
// because validator errors carry no severity. The validator visits children itself.
func Register[T any](v *validator.Validate, check func(T, *Check)) {
	v.RegisterStructValidation(func(sl validator.StructLevel) {
		value, ok := sl.Current().Interface().(T)
		if !ok {
			return
		}
		check(value, &Check{Report: func(violation Violation) {
			if violation.Rule.Severity != Fatal {
				return
			}
			field := violation.Field
			if field == "" {
				field = "-"
			}
			sl.ReportError(value, field, field, violation.Rule.ID, "")
		}})
	}, *new(T))
}
