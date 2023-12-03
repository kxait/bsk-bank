package lib

import (
	"database/sql"
	"fmt"
)

func GetAccountBalance(db *sql.DB, accountId int64) (float64, error) {
	records, err := db.Query("SELECT * FROM account_transaction WHERE source_account_id = ? OR destination_account_id = ?", accountId, accountId)
	if err != nil {
		return 0, err
	}

	transactions, err := ScanTransactions(records)
	if err != nil {
		return 0, err
	}

	balance := 0.0

	for _, transaction := range transactions {
		if transaction.SourceAccountId == accountId {
			balance = balance - transaction.Amount
		}
		if transaction.DestinationAccountId == accountId {
			balance = balance + transaction.Amount
		}

		if balance < 0 {
			return 0, fmt.Errorf("invalid balance! after transaction %d, the balance was %f", transaction.Id, balance)
		}
	}

	return balance, nil
}

func GetAccountsBalances(db *sql.DB) (map[int64]float64, error) {
	balances := make(map[int64]float64, 0)
	records, err := db.Query("SELECT * FROM account_transaction")
	if err != nil {
		return balances, err
	}

	transactions, err := ScanTransactions(records)
	if err != nil {
		return balances, err
	}

	for _, transaction := range transactions {
		sourceAccountBalance, ok := balances[transaction.SourceAccountId]
		if !ok {
			balances[transaction.SourceAccountId] = 0
			sourceAccountBalance = 0
		}
		destinationAccountBalance, ok := balances[transaction.DestinationAccountId]
		if !ok {
			balances[transaction.DestinationAccountId] = 0
			destinationAccountBalance = 0
		}

		if transaction.SourceAccountId != -1 {
			balances[transaction.SourceAccountId] = sourceAccountBalance - transaction.Amount
		}

		if transaction.DestinationAccountId != -1 {
			balances[transaction.DestinationAccountId] = destinationAccountBalance + transaction.Amount
		}

		if transaction.SourceAccountId != -1 && balances[transaction.SourceAccountId] < 0 {
			return make(map[int64]float64, 0), fmt.Errorf("invalid source account balance after transaction %d, balance: %f", transaction.Id, balances[transaction.SourceAccountId])
		}

		if balances[transaction.DestinationAccountId] < 0 {
			return make(map[int64]float64, 0), fmt.Errorf("invalid destination account balance after transaction %d, balance: %f", transaction.Id, balances[transaction.DestinationAccountId])
		}
	}

	return balances, nil
}
