package main

import (
	"fmt"
	"sync"
)

type BankAccount struct {
	balance int
	mutex   sync.Mutex
}

// Deposit money into the account
func (a *BankAccount) Deposit(amount int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.balance += amount
	fmt.Printf("Deposited: %d, New Balance: %d\n", amount, a.balance)
}

// Withdraw money from the account
func (a *BankAccount) Withdraw(amount int) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.balance >= amount {
		a.balance -= amount
		fmt.Printf("Withdrawn: %d, New Balance: %d\n", amount, a.balance)
		return true
	} else {
		fmt.Println("Insufficient balance")
		return false
	}
}

func main() {
	account := &BankAccount{balance: 1000}
	var wg sync.WaitGroup

	// Simulating concurrent deposits and withdrawals
	for i := 0; i < 5; i++ {
		wg.Add(2)
		go func() {
			account.Deposit(500)
			wg.Done()
		}()
		go func() {
			account.Withdraw(300)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Printf("Final Balance: %d\n", account.balance)
}
