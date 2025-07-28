package model

// RefreshTokenRequest is used to request a new access token using a refresh token.
// @Description Refresh your session with a valid refresh token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RegisterRequest contains information for user registration.
// @Description User registration with email/phone, name, password and captcha.
type RegisterRequest struct {
	Identifier string `json:"identifier" example:"user@example.com or 09123456789" binding:"required"`      // Email or mobile number
	FullName   string `json:"full_name" example:"علی منصوری" binding:"required"`                            // Full name of the user
	Password   string `json:"password" example:"P@ssw0rd123" binding:"required"`                            // User password
	CaptchaID  string `json:"captcha_id" binding:"required" example:"a1b2c3d4-e5f6-7g8h-9i10-j11k12l13m14"` // Captcha unique ID
	CaptchaAns string `json:"captcha_ans" binding:"required" example:"aBc123"`                              // Captcha answer
}

// LoginRequest is used for user authentication.
// @Description Login with email/phone, password and captcha.
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required" example:"user@example.com or 09123456789"`      // Email or mobile
	Password   string `json:"password" binding:"required" example:"P@ssw0rd123"`                            // Password
	CaptchaID  string `json:"captcha_id" binding:"required" example:"9f1b77d3-cb13-4a72-8a64-28de5f82a5c2"` // Captcha ID
	CaptchaAns string `json:"captcha_ans" binding:"required" example:"aBc123"`                              // Captcha answer
}

// RegisterResponse represents the response after successful registration.
// @Description User registration response containing user ID and a message.
type RegisterResponse struct {
	UserID  string `json:"user_id" example:"d51c9a7e-0e12-4a7f-abea-1e39b5db9e91"`
	Message string `json:"message" example:"Registration successful."`
}

// Verify2FARequest is used to verify a 2FA code for different purposes (register/login).
// @Description Verify two-factor authentication code (2FA) for registration or login.
type Verify2FARequest struct {
	Code       string `json:"code" binding:"required" example:"123456"`                  // The 2FA code
	Identifier string `json:"identifier"  binding:"required" example:"user@example.com"` // Email or mobile
	Purpose    string `json:"purpose" binding:"required" example:"register_2fa"`         // Purpose: "register_2fa" or "login_2fa"
}

// ResendVerificationRequest is used to resend a verification code.
// @Description Resend a verification code for email or phone verification.
type ResendVerificationRequest struct {
	Identifier string `json:"identifier" binding:"required" example:"user@example.com"` // Email or mobile
	Purpose    string `json:"purpose" example:"register_2fa"`                           // Purpose (optional)
}

// TokenResponse represents access and refresh tokens issued by the authentication system.
// @Description Response containing JWT access and refresh tokens.
type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4uLi4="`
}

// LoginResponse is an alias for TokenResponse (returns tokens after successful login).
type LoginResponse = TokenResponse

// ForgotPasswordRequest is used for initiating password reset.
// @Description Request to send a password reset code to email or mobile.
type ForgotPasswordRequest struct {
	Identifier string `json:"identifier" binding:"required" example:"user@example.com"` // Email or mobile
}

// ResetPasswordRequest is used to reset the user's password using a code.
// @Description Reset password with identifier (email/phone), code, and new password.
type ResetPasswordRequest struct {
	Identifier  string `json:"identifier" binding:"required" example:"user@example.com"` // Email or mobile
	Code        string `json:"code" binding:"required" example:"123456"`                 // The reset code sent to user
	NewPassword string `json:"new_password" binding:"required" example:"N3wP@ssw0rd!"`   // New password
}
