package errcode

var (
	Success          = NewError(0, "success")
	Fail             = NewError(10000000, "internal error")
	InvalidParams    = NewError(10000001, "params error")
	Unauthorized     = NewError(10000002, "Unauthorized error")
	NotFound         = NewError(10000003, "NotFound")
	Unknown          = NewError(10000004, "Unknown error")
	DeadlineExceeded = NewError(10000005, "DeadlineExceeded error")
	AccessDenied     = NewError(10000006, "AccessDenied error")
	LimitExceed      = NewError(10000007, "LimitExceed error")
	MethodNotAllowed = NewError(10000008, "MethodNotAllowed error")
)
