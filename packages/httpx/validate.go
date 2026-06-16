package httpx

import (
	"net/http"
	"strings"

	"ecommerce/packages/types"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DecodeAndValidate(r *http.Request, dst any) error {
	if err := DecodeJSON(r, dst); err != nil {
		return err
	}
	if err := validate.Struct(dst); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			errs := make([]string, 0, len(ve))
			for _, fe := range ve {
				errs = append(errs, fe.Field()+": "+fe.Tag())
			}
			return types.NewValidation(strings.Join(errs, "; "))
		}
	}
	return nil
}
