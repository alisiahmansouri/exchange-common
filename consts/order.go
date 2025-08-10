package consts

// =================== عملیات‌های Handler/UseCase سفارش ===================
const (
	OpOrderPlaceOrder = "OrderHandler.PlaceOrder"
	OpOrderGetByID    = "OrderHandler.GetOrderByID"
	OpOrderList       = "OrderHandler.ListOrders"
	OpOrderCancel     = "OrderHandler.CancelOrder"
)

// =================== پیام‌های موفقیت‌آمیز سفارش ===================
const (
	MsgOrderPlacedSuccessfully   = "سفارش با موفقیت ثبت شد"
	MsgOrderCanceledSuccessfully = "سفارش با موفقیت لغو شد"
	MsgOrderFound                = "سفارش با موفقیت دریافت شد"
	MsgOrdersListedSuccessfully  = "لیست سفارشات با موفقیت دریافت شد"
)

// =================== پیام‌های خطای سفارش ===================
const (
	ErrOrderInvalidBody              = "بدنه درخواست سفارش نامعتبر است"
	ErrOrderCreateFailed             = "خطا در ثبت سفارش"
	ErrOrderNotFound                 = "سفارش یافت نشد"
	ErrOrderUnauthorized             = "دسترسی به سفارش غیرمجاز است"
	ErrOrderInsufficientFunds        = "موجودی کافی برای ثبت سفارش وجود ندارد"
	ErrOrderInvalidAmount            = "مقدار سفارش باید بزرگتر از صفر باشد"
	ErrOrderInvalidID                = "شناسه سفارش نامعتبر است"
	ErrOrderCancelFailed             = "خطا در لغو سفارش"
	ErrOrderCannotBeCanceled         = "این سفارش امکان لغو شدن ندارد"
	ErrOrderAmountOutOfRange         = "مقدار سفارش خارج از بازه مجاز است"
	ErrOrderPairNotFoundOrInactive   = "جفت ارز یافت نشد یا فعال نیست"
	ErrOrderWalletNotFoundOrInactive = "کیف پول یافت نشد یا فعال نیست"
	ErrOrderConflict                 = "سفارش متناقض یا تکراری است"
	ErrOrderTooManyRequests          = "تعداد درخواست‌های سفارش بیش از حد مجاز است"
	ErrOrderTimeout                  = "ثبت سفارش به علت محدودیت زمانی انجام نشد"
	ErrOrderInvalidSide              = "جهت سفارش (خرید/فروش) نامعتبر است"
	ErrOrderClientOrderIDTooLong     = "clientOrderID بیش از حد طولانی است"
	ErrOrderInvalidLimitPrice        = "قیمت سفارش limit نامعتبر است"
	ErrOrderPriceNotAllowedForMarket = "ارسال قیمت برای سفارش market مجاز نیست"
	ErrOrderPairIDInvalid            = "شناسه جفت ارز نامعتبر است"
	ErrOrderInputInvalid             = "داده ورودی سفارش نامعتبر است"
	ErrOrderInvalidType              = "نوع سفارش نامعتبر است"
	ErrInternal                      = "خطای داخلی سرور"
)

// =================== کدهای خطای سفارش (Error Code) ===================
const (
	CodeOrderInvalidBody              = "ORDER_INVALID_BODY"
	CodeOrderCreateError              = "ORDER_CREATE_ERROR"
	CodeOrderNotFound                 = "ORDER_NOT_FOUND"
	CodeOrderUnauthorized             = "ORDER_UNAUTHORIZED"
	CodeOrderInsufficientFunds        = "ORDER_INSUFFICIENT_FUNDS"
	CodeOrderInvalidAmount            = "ORDER_INVALID_AMOUNT"
	CodeOrderInvalidID                = "ORDER_INVALID_ID"
	CodeOrderCancelError              = "ORDER_CANCEL_ERROR"
	CodeOrderCannotBeCanceled         = "ORDER_CANNOT_BE_CANCELED"
	CodeOrderAmountOutOfRange         = "ORDER_AMOUNT_OUT_OF_RANGE"
	CodeOrderPairNotFoundOrInactive   = "ORDER_PAIR_NOT_FOUND_OR_INACTIVE"
	CodeOrderWalletNotFoundOrInactive = "ORDER_WALLET_NOT_FOUND_OR_INACTIVE"
	CodeOrderConflict                 = "ORDER_CONFLICT"
	CodeOrderTooManyRequests          = "ORDER_TOO_MANY_REQUESTS"
	CodeOrderTimeout                  = "ORDER_TIMEOUT"
	CodeOrderInvalidSide              = "ORDER_INVALID_SIDE"
	CodeOrderClientOrderIDTooLong     = "ORDER_CLIENT_ORDER_ID_TOO_LONG"
	CodeOrderInvalidLimitPrice        = "ORDER_INVALID_LIMIT_PRICE"
	CodeOrderPriceNotAllowedForMarket = "ORDER_PRICE_NOT_ALLOWED_FOR_MARKET"
	CodeOrderPairIDInvalid            = "ORDER_PAIR_ID_INVALID"
	CodeOrderInputInvalid             = "ORDER_INPUT_INVALID"
	CodeOrderInvalidType              = "ORDER_INVALID_TYPE"
	CodeInternal                      = "INTERNAL_ERROR"
)
