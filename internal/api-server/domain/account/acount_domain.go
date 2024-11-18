package account

import "github.com/lunarianss/Luna/internal/api-server/repo"

type AccountDomain struct {
	AccountRepo repo.AccountRepo
}

func NewAccountDomain(accountRepo repo.AccountRepo) *AccountDomain {
	return &AccountDomain{
		AccountRepo: accountRepo,
	}
}
