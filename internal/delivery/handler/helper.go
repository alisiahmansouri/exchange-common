package handler

import (
	"exchange-common/internal/captcha"
	"exchange-common/internal/consts"
	"exchange-common/internal/richerror"
	"context"
	"github.com/gin-gonic/gin"
	"strings"
)

// ====== Error Delegators ======

//func richerror.HandleWrap(c *gin.Context, op, userMsg, code string, kind richerror.Kind, err error) {
//	richerror.Handle(c, op, userMsg, code, kind, err)
//}
//
//func HandleWrappedRichError(c *gin.Context, op, userMsg, code string, kind richerror.Kind, err error) {
//	richerror.HandleWrap(c, op, userMsg, code, kind, err)
//}

func verifyCaptchaWrapper(c *gin.Context, op string, ctx context.Context, store *captcha.CaptchaStore, captchaID, captchaAns string) bool {
	if err := verifyCaptcha(ctx, store, captchaID, captchaAns, op); err != nil {
		richerror.HandleWrap(c, op,.ErrCaptchaInvalid, .
		CodeInvalidCaptcha, richerror.KindInvalid, err)
		return true
	}
	return false
}

func verifyCaptcha(ctx context.Context, store *captcha.CaptchaStore, captchaID, captchaAns, op string) error {
	if captchaID == "" || captchaAns == "" {
		return richerror.New(op,.ErrCaptchaInvalidBody, .
		CodeCaptchaEmpty, richerror.KindInvalid, nil)
	}

	expected, err := store.Get(ctx, captchaID)
	if err != nil {
		return richerror.New(op,.ErrCaptchaNotFound, .
		CodeCaptchaNotFound, richerror.KindInvalid, err)
	}

	if strings.ToLower(expected) != strings.ToLower(captchaAns) {
		return richerror.New(op,.ErrCaptchaInvalid, .
		CodeCaptchaWrong, richerror.KindUnauthorized, nil)
	}

	_ = store.Delete(ctx, captchaID)
	return nil
}
