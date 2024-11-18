package code

const (
	// ErrEmailCode - 500: Error occurred when email code is incorrect.
	ErrEmailCode int = iota + 110301
	// ErrTokenEmail - 500: Error occurred when email is incorrect.
	ErrTokenEmail
	// ErrTenantAlreadyExist - 500: Error occurred when tenant is already exist.
	ErrTenantAlreadyExist
)
