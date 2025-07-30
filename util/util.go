package util

import (
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func ParseUUID(value string) (uuid.UUID, error) {
	return uuid.Parse(value)
}

func GenerateUUID() string {
	return uuid.NewString()
}

func GenerateVerificationCode() (string, error) {
	const digits = "0123456789"
	length := 6
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		bytes[i] = digits[bytes[i]%byte(len(digits))]
	}
	return string(bytes), nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("رمز عبور باید حداقل ۸ کاراکتر داشته باشد")
	}

	if len(password) > 64 {
		return errors.New("رمز عبور نمی‌تواند بیشتر از ۶۴ کاراکتر باشد")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// ترکیب حداقل دو نوع کاراکتر
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	if count < 2 {
		return errors.New("رمز عبور باید شامل حداقل دو نوع از موارد زیر باشد: حروف بزرگ، حروف کوچک، اعداد، کاراکترهای ویژه")
	}

	return nil
}

func ValidateEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// بررسی طول ایمیل
	if len(email) > 254 {
		return false
	}

	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return false
	}

	// بررسی دامنه‌های غیرفعال
	blacklistedDomains := []string{"example.com", "test.com"}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.ToLower(parts[1])

	for _, d := range blacklistedDomains {
		if domain == d {
			return false
		}
	}

	return true
}

func Validate2FACode(code string) bool {
	if len(code) == 0 || len(code) > 50 {
		return false
	}
	return true
}

func ValidatePhone(phone string) bool {
	re := regexp.MustCompile(`^09\d{9}$`)
	return re.MatchString(phone)
}
func NormalizeEmail(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func NormalizePhone(s string) string {
	return strings.TrimSpace(s)
}

func Get2FAChannelByIdentifier(identifier string) string {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return ""
	}

	// اگر ایمیل بود
	if ValidateEmail(identifier) {
		return "email"
	}

	// اگر موبایل بود
	if ValidatePhone(identifier) {
		return "sms"
	}

	return ""
}

func NowUTC() time.Time { return time.Now().UTC() }
