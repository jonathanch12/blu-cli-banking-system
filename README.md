# BLU CLI Banking System

A file-based banking system built in Go as a command-line application. It manages user accounts, transfers funds between them, records savings deposits, and compounds interest on those deposits over time. All state is persisted as plain text files, so no database is required.

## Features

- **Create accounts** with a starting balance.
- **Transfer funds** between two existing accounts, recorded in both users' files.
- **Add savings deposits** that are tracked separately from the main account balance.
- **Accrue interest** on all deposits at a fixed rate of 1% per minute, compounded.

### System Rules

- Interest is 1% (`0.01`) per minute, compounded based on the whole number of minutes elapsed since a deposit's last update.
- All account, deposit, and history data is stored in plain text files under `filebase/`.
- Deposits are **not** deducted from the main account balance; they are tracked independently.
- Deposits younger than one minute are left untouched when interest is accrued.

## Directory Structure

```
blu-cli-banking-system/
├── go.mod                  # Go module definition (module: blu)
├── main.go                 # CLI entry point: parses commands and arguments
├── case.txt                # Original project brief / requirements
├── services/               # Core business logic
│   ├── account.go          # CreateAccount, Transfer, balance reading
│   ├── deposit.go          # AddDeposit, AccrueInterest
│   └── helper.go           # filePathFor: builds file paths under filebase/
└── filebase/               # Persistent data store (plain text)
    ├── accounts/           # <name>.txt — main balance + transaction log
    ├── deposits/           # <name>.txt — one line per deposit
    └── history/            # <name>.txt — chronological action log
```

### File Formats

**`filebase/accounts/<name>.txt`** — first line is the current balance, followed by a transaction log:

```
500.00
CREATE|1000.00|2026-09-25T18:12:50+07:00
TRANSFER|500.00|TestUserB|2026-09-25T18:15:36+07:00
```

**`filebase/deposits/<name>.txt`** — one deposit per line, formatted as `<id>|<amount>|<timestamp>`:

```
1|12072.84|2026-09-25T23:33:59+07:00
2|7243.71|2026-09-25T23:33:59+07:00
```

**`filebase/history/<name>.txt`** — chronological log of every action affecting the user:

```
CREATE|1000.00|2026-09-25T18:12:50+07:00
DEPOSIT|1|500.00|2026-09-25T18:13:39+07:00
TRANSFER|500.00|TestUserB|2026-09-25T18:15:36+07:00
INTEREST|1|12072.84|2026-09-25T23:33:59+07:00
```

## Prerequisites

- [Go](https://go.dev/dl/) 1.27.1 or later (per `go.mod`).
- The `filebase/accounts`, `filebase/deposits`, and `filebase/history` directories must exist before running commands, since the program writes into them.

## Usage

Run commands from the project root using `go run main.go <command> [arguments]`.

### 1. Create an account

```
go run main.go create_account <name> <amount>
```

Creates an account file with the starting balance and initializes matching deposit and history files.

```
go run main.go create_account Alice 1000.00
```

### 2. Transfer funds

```
go run main.go transfer <from> <to> <amount>
```

Moves funds from the sender to the receiver. Both accounts must exist, the amount must be positive, and the sender must have a sufficient balance.

```
go run main.go transfer Rudolph Alice 500.00
```

### 3. Add a deposit

```
go run main.go add_deposit <name> <amount>
```

Appends a deposit entry (with an auto-incrementing id and timestamp) to the user's deposit file. Does not affect the main account balance.

```
go run main.go add_deposit Alice 500.00
```

### 4. Accrue interest

```
go run main.go accrue_interest
```

Scans every deposit file and compounds interest at 1% per minute based on the minutes elapsed since each deposit's last update. Updated deposits get an `INTEREST` entry appended to the user's history.

## Example Session

```
go run main.go create_account Alice 1000.00
go run main.go create_account Rudolph 2000.00
go run main.go add_deposit Alice 500.00
go run main.go transfer Rudolph Alice 500.00
# wait a few minutes...
go run main.go accrue_interest
```

After this sequence:

- Alice's account is created with a balance of 1000.00, Rudolph's with 2000.00.
- Alice records a 500.00 deposit (main balance unchanged).
- Rudolph transfers 500.00 to Alice, leaving Rudolph at 1500.00 and Alice at 1500.00.
- Running `accrue_interest` compounds Alice's 500.00 deposit at 1% per minute for the elapsed minutes, with no change to any main balance.

## Building a Binary

To compile a standalone executable instead of using `go run`:

```
go build -o blu .
```

Then invoke it directly, for example:

```
./blu create_account Alice 1000.00
```
