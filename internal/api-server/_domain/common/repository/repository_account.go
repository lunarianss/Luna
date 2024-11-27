package account

type BaseAccount interface {
	GetAccountType() string
	GetAccountID() string
}
