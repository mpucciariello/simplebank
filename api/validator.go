package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/micaelapucciariello/simplebank/utils"
)

// build a validator to test different currencies in a more scalable way

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	currency, ok := fieldLevel.Field().Interface().(string)
	if ok {
		utils.IsSupported(currency)
	}
	return false
}
