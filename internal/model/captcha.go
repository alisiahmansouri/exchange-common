package model

type CaptchaResponse struct {
	RequestID string `json:"request_id"`
	Captcha   string `json:"captcha"`
}

type VerifyCaptchaRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	UserInput string `json:"user_input" binding:"required"`
}
