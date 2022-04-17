package accounts

import "errors"

type Account struct {
	owner   string
	balance int
}

// NewAccount creates  new account
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

// Deposit x amount on your account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

//Balance of your account
func (a Account) Balance() int {
	return a.balance
}

// Withdraw x amount
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errors.New("넌 거지야 돈 못뽑아!")
	}
	a.balance -= amount
	return nil
}

// ChangeOwner new owner
func (a *Account) ChangeOwner(owner string) {
	a.owner = owner
}

func (a Account) Owner() string {
	return a.owner
}
