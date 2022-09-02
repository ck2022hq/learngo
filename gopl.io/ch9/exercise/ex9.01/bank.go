// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

// Package bank provides a concurrency-safe bank with one account.
package bank

type withdraws struct {
	amount  int
	success chan bool
}

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance
var withdraw = make(chan withdraws)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) bool {
	success := make(chan bool)
	withdraw <- withdraws{amount, success}
	return <-success
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case wd := <-withdraw:
			if wd.amount <= balance {
				balance -= wd.amount
				wd.success <- true
			} else {
				wd.success <- false
			}
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-
