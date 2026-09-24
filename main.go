package main

import (
	"blu/services"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No arguments detected.")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create_account":
		if len(os.Args) < 4 {
			fmt.Println("Missing arguments for account creation.")
			os.Exit(1)
		}
		name := os.Args[2]
		amount, err := strconv.ParseFloat(os.Args[3], 64)

		if err != nil {
			fmt.Printf("Error detected: %v", err)
			os.Exit(1)
		}

		accError := services.CreateAccount(name, amount)
		if accError != nil {
			fmt.Fprintln(os.Stderr, accError)
			os.Exit(1)
		}

		fmt.Printf("Account creation for %s is successful.\n", name)
	case "transfer":
		if len(os.Args) < 5 {
			fmt.Println("Missing arguments for performing transfer.")
			os.Exit(1)
		}

		sender := os.Args[2]
		receiver := os.Args[3]
		amount, err := strconv.ParseFloat(os.Args[4], 64)

		if err != nil {
			fmt.Printf("Error detected: %v", err)
			os.Exit(1)
		}

		transferError := services.Transfer(sender, receiver, amount)
		if transferError != nil {
			fmt.Fprintln(os.Stderr, transferError)
			os.Exit(1)
		}

		fmt.Printf("%.2f is successfully transfered from %s to %s.\n", amount, sender, receiver)
		//	case "add_deposit" :
		//	case "accrue_interest":
	default:
		fmt.Fprintln(os.Stderr, "Command not recognized.")
		os.Exit(1)
	}
}
