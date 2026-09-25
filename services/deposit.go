package services

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func AddDeposit(name string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount cannot be less than zero, inputted amount: %.2f", amount)
	}

	userDeposit := filePathFor("deposits", name)
	if _, err := os.Stat(userDeposit); err != nil {
		return fmt.Errorf("deposit account for %q does not exist", name)
	}

	userHistory := filePathFor("history", name)
	if _, err := os.Stat(userHistory); err != nil {
		return fmt.Errorf("history account for %q does not exist", name)
	}

	nextId := 1
	userDepositData, err := os.ReadFile(userDeposit)
	if err != nil {
		return fmt.Errorf("error reading user's deposit account: %w", err)
	}
	trimmed := strings.TrimSpace(string(userDepositData))
	if trimmed != "" {
		lines := strings.Split(trimmed, "\n")
		lastLine := lines[len(lines)-1]
		data := strings.Split(lastLine, "|")
		lastId, err := strconv.Atoi(strings.TrimSpace(data[0]))
		if err != nil {
			return fmt.Errorf("error parsing latest deposit id: %w", err)
		}
		nextId = lastId + 1
	}

	ts := time.Now().Format(time.RFC3339)

	fd, err := os.OpenFile(userDeposit, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening user deposit: %w", err)
	}
	defer fd.Close()
	if _, err := fmt.Fprintf(fd, "%d|%.2f|%s\n", nextId, amount, ts); err != nil {
		return fmt.Errorf("error writing deposit for user %q: %w", name, err)
	}

	fh, err := os.OpenFile(userHistory, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening user history: %w", err)
	}
	defer fh.Close()
	if _, err := fmt.Fprintf(fh, "DEPOSIT|%d|%.2f|%s\n", nextId, amount, ts); err != nil {
		return fmt.Errorf("error writing user history: %w", err)
	}

	return nil
}
