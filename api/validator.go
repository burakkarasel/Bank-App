package api

import (
	"github.com/burakkarasel/Bank-App/util"
	"github.com/go-playground/validator/v10"
)

// validCurrency is a custom validator that checks if a given currency is valid or not
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is supported
		return util.IsSupportedCurrency(currency)
	}

	return false
}
