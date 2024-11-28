package repository

type BaseAccount interface {
	GetAccountType() string
	GetAccountID() string
}
