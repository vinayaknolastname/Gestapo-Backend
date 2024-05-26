package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/akmal4410/gestapo/pkg/utils"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateBody(body io.Reader, data any) error {
	RegisterValidator()
	if body != nil {
		if err := json.NewDecoder(body).Decode(&data); err != nil {
			return err
		}
	}

	err := validate.Struct(data)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		var errorMessage strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			errorMessage.WriteString(fmt.Sprintf("Field '%s' failed validation. ", fieldName))
		}

		return fmt.Errorf("validation errors: %s", errorMessage.String())
	}
	return nil

}

func RegisterValidator() {
	err := validate.RegisterValidation("user_type", validateUserType)
	if err != nil {
		fmt.Println("Error registering user_type:", err.Error())
	}
	err = validate.RegisterValidation("signup_action", validateSignupAction)
	if err != nil {
		fmt.Println("Error registering signup_action:", err.Error())
	}
	err = validate.RegisterValidation("sso_action", validateSSOAction)
	if err != nil {
		fmt.Println("Error registering sso_action:", err.Error())
	}
	err = validate.RegisterValidation("gender", validateGender)
	if err != nil {
		fmt.Println("Error registering gender:", err.Error())
	}
	err = validate.RegisterValidation("validate_date", validateDate)
	if err != nil {
		fmt.Println("Error registering validate_date:", err.Error())
	}
	err = validate.RegisterValidation("percentage", validatePercentage)
	if err != nil {
		fmt.Println("Error registering percentage:", err.Error())
	}
	err = validate.RegisterValidation("wishlist_action", validateAddRemoveWishlistAction)
	if err != nil {
		fmt.Println("Error registering wishlist_action:", err.Error())
	}
	err = validate.RegisterValidation("payment_mode", validatePaymentMode)
	if err != nil {
		fmt.Println("Error registering payment_mode:", err.Error())
	}
	err = validate.RegisterValidation("order_type", validateOrderType)
	if err != nil {
		fmt.Println("Error registering order_type:", err.Error())
	}
}

var validateUserType validator.Func = func(fl validator.FieldLevel) bool {
	if userType, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedUsers(userType)
	}
	return false
}

var validateSignupAction validator.Func = func(fl validator.FieldLevel) bool {
	if signupAction, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedSignupAction(signupAction)
	}
	return false
}

var validateSSOAction validator.Func = func(fl validator.FieldLevel) bool {
	if signupAction, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedSSOAction(signupAction)
	}
	return false
}

var validateGender validator.Func = func(fl validator.FieldLevel) bool {
	if signupAction, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedGender(signupAction)
	}
	return false
}

var validateDate validator.Func = func(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

var validatePercentage validator.Func = func(fl validator.FieldLevel) bool {
	if perentage, ok := fl.Field().Interface().(float64); ok {
		return utils.IsSupportedPercentage(perentage)
	}
	return false
}

var validateAddRemoveWishlistAction validator.Func = func(fl validator.FieldLevel) bool {
	if wishlistAction, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedAddRemoveWishlistAction(wishlistAction)
	}
	return false
}

var validatePaymentMode validator.Func = func(fl validator.FieldLevel) bool {
	if wishlistAction, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedPaymentMode(wishlistAction)
	}
	return false
}

var validateOrderType validator.Func = func(fl validator.FieldLevel) bool {
	if wishlistAction, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedOrderTypeMode(wishlistAction)
	}
	return false
}
