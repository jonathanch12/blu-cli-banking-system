package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func filePathFor(kind, name string) string {
	return filepath.Join("filebase", kind, name+".txt")
}

func readBalance(name string) (float64, error) {
	path := filePathFor("accounts", name)
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("error reading account %q: %w", name, err)
	}

	lines := strings.Split(string(data), "\n")

	balance, err := strconv.ParseFloat(strings.TrimSpace(lines[0]), 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing balance for %q: %w", name, err)
	}

	return balance, nil
}

func CreateAccount(name string, amount float64) error {
	filePath := filePathFor("accounts", name)
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("account %q already exists", name)
	}

	ts := time.Now().Format(time.RFC3339)

	content := fmt.Sprintf("%.2f\nCREATE,%.2f,%s\n", amount, amount, ts)
	err := os.WriteFile(filePath, []byte(content), 0644)

	if err != nil {
		return fmt.Errorf("account creation failed, error: %w", err)
	}

	depositPath := filePathFor("deposits", name)
	err1 := os.WriteFile(depositPath, []byte(""), 0644)

	if err1 != nil {
		return fmt.Errorf("deposit account creation failed, error: %w", err1)
	}

	historyPath := filePathFor("history", name)
	historyLine := fmt.Sprintf("CREATE,%.2f,%s\n", amount, ts)
	err2 := os.WriteFile(historyPath, []byte(historyLine), 0644)

	if err2 != nil {
		return fmt.Errorf("failed to save account creation log, error: %w", err2)
	}
	return nil
}

func Transfer(sender string, receiver string, amount float64) error {
	senderPath := filePathFor("accounts", sender)
	if _, err := os.Stat(senderPath); err != nil {
		return fmt.Errorf("sender account %q not found", sender)
	}

	receiverPath := filePathFor("accounts", receiver)
	if _, err := os.Stat(receiverPath); err != nil {
		return fmt.Errorf("receiver account %q not found", receiver)
	}

	if amount <= 0 {
		return fmt.Errorf("amount cannot be less than zero, inputted amount: %.2f", amount)
	}

	senderBalance, err := readBalance(sender)
	if err != nil {
		return err
	}

	if senderBalance < amount {
		return fmt.Errorf("sender %q does not have sufficient amount to transfer", sender)
	}

	receiverBalance, err := readBalance(receiver)
	if err != nil {
		return err
	}

	newSenderBalance := senderBalance - amount
	newReceiverBalance := receiverBalance + amount

	senderData, err := os.ReadFile(senderPath)
	if err != nil {
		return fmt.Errorf("error reading sender's account: %w", err)
	}
	senderDataLines := strings.Split(strings.TrimSpace(string(senderData)), "\n")
	senderLogLines := senderDataLines[1:]

	receiverData, err := os.ReadFile(receiverPath)
	if err != nil {
		return fmt.Errorf("error reading receiver's account: %w", err)
	}
	receiverDataLines := strings.Split(strings.TrimSpace(string(receiverData)), "\n")
	receiverLogLines := receiverDataLines[1:]

	ts := time.Now().Format(time.RFC3339)

	oldSenderLog := strings.Join(senderLogLines, "\n")
	senderNewLog := fmt.Sprintf("%.2f\n%s\nTRANSFER,%.2f,%s,%s\n", newSenderBalance, oldSenderLog, amount, receiver, ts)

	oldReceiverLog := strings.Join(receiverLogLines, "\n")
	receiverNewLog := fmt.Sprintf("%.2f\n%s\nRECEIVE,%.2f,%s,%s\n", newReceiverBalance, oldReceiverLog, amount, sender, ts)

	if err := os.WriteFile(senderPath, []byte(senderNewLog), 0644); err != nil {
		return fmt.Errorf("error updating sender account: %w", err)
	}
	if err := os.WriteFile(receiverPath, []byte(receiverNewLog), 0644); err != nil {
		return fmt.Errorf("error updating receiver account: %w", err)
	}

	senderHist := filePathFor("history", sender)
	fs, err := os.OpenFile(senderHist, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening sender history: %w", err)
	}
	defer fs.Close()
	if _, err := fmt.Fprintf(fs, "TRANSFER,%.2f,%s,%s\n", amount, receiver, ts); err != nil {
		return fmt.Errorf("error writing sender history: %w", err)
	}

	receiverHist := filePathFor("history", receiver)
	fr, err := os.OpenFile(receiverHist, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening receiver history: %w", err)
	}
	defer fr.Close()
	if _, err := fmt.Fprintf(fr, "RECEIVE,%.2f,%s,%s\n", amount, sender, ts); err != nil {
		return fmt.Errorf("error writing receiver history: %w", err)
	}

	return nil
}
