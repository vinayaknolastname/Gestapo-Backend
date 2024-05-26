package entity

type SendOTPReq struct {
	Email  string `json:"email" validate:"omitempty"`
	Phone  string `json:"phone" validate:"omitempty,len=10,numeric"`
	Action string `json:"action" validate:"signup_action"`
}

type SignupReq struct {
	Email    string `json:"email" validate:"omitempty"`
	Phone    string `json:"phone" validate:"omitempty,len=10,numeric"`
	FullName string `json:"full_name" validate:"required,min=4,max=20"`
	UserName string `json:"user_name" validate:"required,min=4,max=12"`
	UserType string `json:"user_type" validate:"user_type"`
	Code     string `json:"code" validate:"required,min=6,max=6"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type SsoReq struct {
	UserType string `json:"user_type" validate:"user_type"`
	Action   string `json:"action" validate:"sso_action"`
}

type LoginReq struct {
	UserName string `json:"user_name" validate:"required,min=4,max=12"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type ForgotPassword struct {
	Email    string `json:"email" validate:"omitempty"`
	Phone    string `json:"phone" validate:"omitempty,len=10,numeric"`
	Code     string `json:"code" validate:"required,min=6,max=6"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}
