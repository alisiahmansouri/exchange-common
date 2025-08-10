package consts

const (
	CodeUnauthorized   = "UNAUTHORIZED"
	CodeForbidden      = "FORBIDDEN"
	CodeNotFound       = "NOT_FOUND"
	CodeInvalidRequest = "INVALID_REQUEST"
	CodeInternalError  = "INTERNAL_ERROR"
)

// ──────────────────────────────
// عملیات (برای لاگ و richerror)
// ──────────────────────────────
const (
	OpWalletCreate        = "WalletHandler.CreateWallet"
	OpWalletDeposit       = "WalletHandler.Deposit"
	OpWalletWithdraw      = "WalletHandler.Withdraw"
	OpWalletFreeze        = "WalletHandler.Freeze"
	OpWalletUnfreeze      = "WalletHandler.Unfreeze"
	OpWalletDeductFrozen  = "WalletHandler.DeductFrozen"
	OpWalletTransfer      = "WalletHandler.TransferInternal"
	OpWalletChangeStatus  = "WalletHandler.ChangeWalletStatus"
	OpWalletSummary       = "WalletHandler.GetWalletsSummary"
	OpWalletHistory       = "WalletHandler.GetWalletHistory"
	OpWalletBulkOperation = "WalletHandler.BulkWalletOperation"
	OpWalletGetByID       = "WalletHandler.GetWalletByID" // ✅ اضافه شد
)

// ──────────────────────────────
// پیام خطا (برای کاربر)
// ──────────────────────────────
const (
	ErrWalletInvalidBody         = "بدنه درخواست کیف پول نامعتبر است"
	ErrWalletAlreadyExists       = "کیف پول قبلاً وجود دارد"
	ErrWalletCreateFailed        = "خطا در ایجاد کیف پول"
	ErrWalletNotFound            = "کیف پول یافت نشد"
	ErrWalletUnauthorized        = "دسترسی به کیف پول غیرمجاز است"
	ErrWalletInsufficientFunds   = "موجودی کیف پول کافی نیست"
	ErrWalletFrozenInsufficient  = "موجودی فریز شده کافی نیست"
	ErrWalletInvalidAmount       = "مبلغ باید بزرگتر از صفر باشد"
	ErrWalletInvalidID           = "شناسه کیف پول نامعتبر است"
	ErrWalletInvalidCurrencyID   = "شناسه ارز نامعتبر است"
	ErrWalletInactive            = "کیف پول غیر فعال است"
	ErrWalletInvalidStatus       = "وضعیت کیف پول نامعتبر است"
	ErrWalletDepositFailed       = "خطا در انجام عملیات واریز به کیف پول"
	ErrWalletWithdrawFailed      = "خطا در برداشت از کیف پول"
	ErrWalletFreezeFailed        = "خطا در فریز کردن کیف پول"
	ErrWalletUnfreezeFailed      = "خطا در آزادسازی موجودی فریز شده"
	ErrWalletDeductFrozenFailed  = "خطا در کسر از موجودی فریز شده"
	ErrWalletTransferFailed      = "خطا در انتقال داخلی بین کیف پول‌ها"
	ErrWalletChangeStatusFailed  = "خطا در تغییر وضعیت کیف پول"
	ErrWalletUserInvalidID       = "شناسه کاربر نامعتبر است"
	ErrWalletListFailed          = "خطا در دریافت لیست کیف پول‌ها"
	ErrWalletSummaryFailed       = "خطا در دریافت خلاصه موجودی‌ها"
	ErrWalletHistoryFailed       = "خطا در دریافت تاریخچه کیف پول"
	ErrWalletBulkOperationFailed = "خطا در انجام عملیات گروهی کیف پول"
	ErrWalletFetchFailed         = "خطا در دریافت اطلاعات کیف پول" // ✅ اضافه شد
)

// ──────────────────────────────
// کد خطا (برای فرانت‌اند/کلاینت)
// ──────────────────────────────
const (
	CodeInvalidWalletID     = "INVALID_WALLET_ID"
	CodeInvalidCurrencyID   = "INVALID_CURRENCY_ID"
	CodeInvalidAmount       = "INVALID_AMOUNT"
	CodeInvalidWalletStatus = "INVALID_WALLET_STATUS"

	CodeWalletNotFound = "WALLET_NOT_FOUND"
	CodeWalletInactive = "WALLET_INACTIVE"

	CodeDepositError       = "DEPOSIT_ERROR"
	CodeWithdrawError      = "WITHDRAW_ERROR"
	CodeFreezeError        = "FREEZE_ERROR"
	CodeUnfreezeError      = "UNFREEZE_ERROR"
	CodeDeductFrozenError  = "DEDUCT_FROZEN_ERROR"
	CodeTransferError      = "TRANSFER_ERROR"
	CodeChangeStatusError  = "CHANGE_STATUS_ERROR"
	CodeWalletListError    = "WALLET_LIST_ERROR"
	CodeWalletSummaryError = "WALLET_SUMMARY_ERROR"
	CodeWalletHistoryError = "WALLET_HISTORY_ERROR"
	CodeBulkOperationError = "BULK_OPERATION_ERROR"
	CodeFetchError         = "FETCH_ERROR" // ✅ اضافه شد
)

// ──────────────────────────────
// پیام موفقیت
// ──────────────────────────────
const (
	MsgWalletCreateSuccess  = "کیف پول با موفقیت ایجاد شد"
	MsgDepositSuccess       = "واریز با موفقیت انجام شد"
	MsgWithdrawSuccess      = "برداشت با موفقیت انجام شد"
	MsgFreezeSuccess        = "موجودی با موفقیت فریز شد"
	MsgUnfreezeSuccess      = "موجودی فریز شده با موفقیت آزاد شد"
	MsgDeductFrozenSuccess  = "کسر از موجودی فریز شده با موفقیت انجام شد"
	MsgTransferSuccess      = "انتقال داخلی با موفقیت انجام شد"
	MsgChangeStatusSuccess  = "وضعیت کیف پول با موفقیت تغییر کرد"
	MsgBulkOperationSuccess = "عملیات گروهی کیف پول با موفقیت انجام شد"
	MsgWalletFetchSuccess   = "اطلاعات کیف پول با موفقیت دریافت شد" // ✅ اضافه شد
)
