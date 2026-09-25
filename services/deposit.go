package services

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
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

// AccrueInterest scans every deposit file under filebase/deposits and compounds
// interest at 1% per minute on each deposit, based on the whole number of minutes
// elapsed since that deposit's last-recorded timestamp. Deposits younger than one
// minute are left untouched (amount and timestamp unchanged). For each deposit that
// is updated, an INTEREST line is appended to the user's history file.
func AccrueInterest() error {
	const ratePerMinute = 0.01

	depositsDir := filepath.Join("filebase", "deposits")
	entries, err := os.ReadDir(depositsDir)
	if err != nil {
		return fmt.Errorf("error reading deposits directory: %w", err)
	}

	now := time.Now()

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".txt")
		depositPath := filePathFor("deposits", name)

		data, err := os.ReadFile(depositPath)
		if err != nil {
			return fmt.Errorf("error reading deposit file for %q: %w", name, err)
		}

		trimmed := strings.TrimSpace(string(data))
		if trimmed == "" {
			// No deposits for this user; nothing to accrue.
			continue
		}

		lines := strings.Split(trimmed, "\n")
		newLines := make([]string, 0, len(lines))

		// Collect the history entries to append only after every line parses
		// successfully, so a malformed line aborts before any file is touched.
		type interestLog struct {
			id     int
			amount float64
			ts     string
		}
		var logs []interestLog

		for _, line := range lines {
			fields := strings.Split(line, "|")
			if len(fields) != 3 {
				return fmt.Errorf("malformed deposit line for %q: %q", name, line)
			}

			id, err := strconv.Atoi(strings.TrimSpace(fields[0]))
			if err != nil {
				return fmt.Errorf("error parsing deposit id for %q in line %q: %w", name, line, err)
			}

			amount, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
			if err != nil {
				return fmt.Errorf("error parsing deposit amount for %q in line %q: %w", name, line, err)
			}

			lastTs, err := time.Parse(time.RFC3339, strings.TrimSpace(fields[2]))
			if err != nil {
				return fmt.Errorf("error parsing deposit timestamp for %q in line %q: %w", name, line, err)
			}

			minutes := int(now.Sub(lastTs).Minutes())
			if minutes <= 0 {
				// Younger than a minute (or clock skew): leave the line untouched.
				newLines = append(newLines, line)
				continue
			}

			newAmount := amount * math.Pow(1+ratePerMinute, float64(minutes))
			newTs := now.Format(time.RFC3339)

			newLines = append(newLines, fmt.Sprintf("%d|%.2f|%s", id, newAmount, newTs))
			logs = append(logs, interestLog{id: id, amount: newAmount, ts: newTs})
		}

		// Rewrite the deposit file with the recomputed lines.
		newContent := strings.Join(newLines, "\n") + "\n"
		if err := os.WriteFile(depositPath, []byte(newContent), 0644); err != nil {
			return fmt.Errorf("error writing deposit file for %q: %w", name, err)
		}

		// Append INTEREST entries to history for the deposits that changed.
		if len(logs) == 0 {
			continue
		}

		historyPath := filePathFor("history", name)
		fh, err := os.OpenFile(historyPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("error opening history file for %q: %w", name, err)
		}
		for _, l := range logs {
			if _, err := fmt.Fprintf(fh, "INTEREST|%d|%.2f|%s\n", l.id, l.amount, l.ts); err != nil {
				fh.Close()
				return fmt.Errorf("error writing interest history for %q: %w", name, err)
			}
		}
		fh.Close()
	}

	return nil
}
