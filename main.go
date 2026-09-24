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

		//	case "add_deposit" :
		//	case "accrue_interest":
	default:
		fmt.Fprintln(os.Stderr, "Command not recognized.")
		os.Exit(1)
	}
}
