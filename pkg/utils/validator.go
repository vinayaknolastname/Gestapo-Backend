package utils

// constants for all supported SignupAction
const (
	SIGN_UP         string = "sign-up"
	FORGOT_PASSWORD string = "forgot-password"
)

// constants for all supported SSO Action
const (
	SSO_ANDROID string = "sso-android"
	SSO_IOS     string = "sso-ios"
)

// constants for all supported USER_TYPE
const (
	USER     = "USER"
	MERCHANT = "MERCHANT"
	ADMIN    = "ADMIN"
)

// constants for all supported Gender
const (
	MALE   = "Male"
	FEMALE = "Female"
)

// constants for all supported SignupAction
const (
	ADD_WISHLIST    string = "add-wishlist"
	REMOVE_WISHLIST string = "remove-wishlist"
)

// constants for all supported Payment Mode
const (
	COD   string = "COD"
	OTHER string = "OTHER"
)

// IsSupportedSignupAction returns true if the SignupAction is supported
func IsSupportedSignupAction(action string) bool {
	switch action {
	case SIGN_UP, FORGOT_PASSWORD:
		return true
	}
	return false
}

// IsSupportedSSOAction returns true if the SSOAction is supported
func IsSupportedSSOAction(action string) bool {
	switch action {
	case SSO_ANDROID, SSO_IOS:
		return true
	}
	return false
}

// IsSupportedUsers returns true if the USER_TYPE is supported
func IsSupportedUsers(usertType string) bool {
	switch usertType {
	case USER, MERCHANT, ADMIN:
		return true
	}
	return false
}

// IsSupportedUsers returns true if the Gender is supported
func IsSupportedGender(gender string) bool {
	switch gender {
	case MALE, FEMALE:
		return true
	}
	return false
}

// IsSupportedPercentage returns true if the percentage is supported
func IsSupportedPercentage(percentage float64) bool {
	if percentage > 1 && percentage < 100 {
		return true
	}
	return false
}

func IsValidPassword(password string) bool {
	if len(password) > 6 && len(password) < 100 {
		return true
	}
	return false
}

func IsValidCode(code string) bool {
	return len(code) == 6
}

// IsSupportedAddRemoveWishlistAction returns true if the AddRemove action is supported
func IsSupportedAddRemoveWishlistAction(action string) bool {
	switch action {
	case ADD_WISHLIST, REMOVE_WISHLIST:
		return true
	}
	return false
}

// IsSupportedPaymentMode returns true if the Payment mode is supported
func IsSupportedPaymentMode(mode string) bool {
	switch mode {
	case COD, OTHER:
		return true
	}
	return false
}

// IsSupportedOrderTypeMode returns true if the Order mode is supported
func IsSupportedOrderTypeMode(mode string) bool {
	switch mode {
	case OrderActive, OrderCompleted, OrderCancelled:
		return true
	}
	return false
}
