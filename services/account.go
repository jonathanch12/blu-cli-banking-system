package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func filePathFor(kind, name string) string {
	return filepath.Join("filebase", kind, name+".txt")
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
