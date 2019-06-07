package models

type Account struct {
	Key     string
	Secret  string
	EventID uint32
	Blocked bool
}

func (account *Account) CheckHash(hash string) bool {
	// return !account.Blocked && hash == tools.MD5(account.Key + account.Secret)
	return false
}

