package main

import (
	"accounts"
	"dictionary"
	"fmt"
	"log"
)

func bank() {
	account := accounts.NewAccount("juno")

	account.Deposit(11)

	err := account.Withdraw(10)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("total balance : %d\n", account.GetBalance())
}

func dict() {
	dictionary := dictionary.Dictionary{"hello": "안녕"}
	dictionary.Add("hello2", "안녕2")

	value, err2 := dictionary.Search("hello2")
	if err2 != nil {
		fmt.Println(err2)

		fmt.Println(dictionary)
	} else {
		fmt.Println(value)
	}

}

func main() {
	dict()
}
