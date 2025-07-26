package model

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type RegisterRequest struct {
	Identifier string `json:"identifier" example:"user@example.com|9209201595" binding:"required"`
	FullName   string `json:"full_name" example:"علی منصوری" binding:"required"`
	Password   string `json:"password" example:"12345678" binding:"required"`
	CaptchaID  string `json:"captcha_id" binding:"required"`
	CaptchaAns string `json:"captcha_ans" binding:"required"`
}
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required" example:"user@example.com or 09123456789"`
	Password   string `json:"password" binding:"required" example:"P@ssw0rd123"`
	CaptchaID  string `json:"captcha_id" binding:"required" example:"9f1b77d3-cb13-4a72-8a64-28de5f82a5c2"`
	CaptchaAns string `json:"captcha_ans" binding:"required" example:"aBc123"`
}

type RegisterResponse struct {
	UserID  string `json:"user_id" example:"d51c9a7e-0e12-4a7f-abea-1e39b5db9e91"`
	Message string `json:"message" example:"Message"`
}

type Verify2FARequest struct {
	Code       string `json:"code" binding:"required"`
	Identifier string `json:"identifier"  binding:"required"`
	Purpose    string `json:"purpose" binding:"required"`
}

type ResendVerificationRequest struct {
	Identifier string `json:"identifier" binding:"required"` // ایمیل یا موبایل
	Purpose    string `json:"purpose"`                       // اختیاری: مثلا "register_2fa" یا "login_2fa"
}

type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4uLi4="`
}

type LoginResponse = TokenResponse

type ForgotPasswordRequest struct {
	Identifier string `json:"identifier" binding:"required"` // ایمیل یا موبایل
}

type ResetPasswordRequest struct {
	Identifier  string `json:"identifier" binding:"required"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type SendPhoneVerificationRequest struct {
	Phone string `json:"phone" binding:"required"`
}
type VerifyPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type ResendEmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResendPhoneVerificationRequest struct {
	Phone string `json:"phone" binding:"required"`
}
