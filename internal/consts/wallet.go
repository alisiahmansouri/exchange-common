package consts

const (
	OpWalletCreate   = "WalletHandler.CreateWallet"
	OpWalletDeposit  = "WalletHandler.Deposit"
	OpWalletWithdraw = "WalletHandler.Withdraw"
)

const (
	ErrWalletInvalidBody       = "بدنه درخواست کیف پول نامعتبر است"
	ErrWalletAlreadyExists     = "کیف پول قبلاً وجود دارد"
	ErrWalletCreateFailed      = "خطا در ایجاد کیف پول"
	ErrWalletNotFound          = "کیف پول یافت نشد"
	ErrWalletUnauthorized      = "دسترسی به کیف پول غیرمجاز است"
	ErrWalletInsufficientFunds = "موجودی کیف پول کافی نیست"
	ErrWalletInvalidAmount     = "مبلغ باید بزرگتر از صفر باشد"
	ErrWalletInvalidID         = "شناسه کیف پول نامعتبر است"
	ErrWalletDepositFailed     = "خطا در انجام عملیات واریز به کیف پول"
	ErrWalletWithdrawFailed    = "خطا در برداشت از کیف پول"
	ErrWalletUserInvalidID     = "شناسه کاربر نامعتبر است"
	ErrWalletListFailed        = "خطا در دریافت لیست کیف پول‌ها"
)

const (
	CodeInvalidCurrencyID = "INVALID_CURRENCY_ID"
	CodeInvalidWalletID   = "INVALID_WALLET_ID"
	CodeDepositError      = "DEPOSIT_ERROR"
	CodeWithdrawError     = "WITHDRAW_ERROR"
	CodeWalletListError   = "WALLET_LIST_ERROR"
)

const (
	MsgDepositSuccess  = "واریز با موفقیت انجام شد"
	MsgWithdrawSuccess = "برداشت با موفقیت انجام شد"
)
