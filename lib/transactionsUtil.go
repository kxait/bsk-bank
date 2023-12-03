package lib

import (
	"database/sql"
	"fmt"
)

func CreateTransaction(db *sql.DB, sourceAccountId int64, destinationAccountId int64, amount float64) (*Transaction, error) {
	records, err := db.Query("SELECT * FROM account WHERE id = ?", sourceAccountId)
	if err != nil {
		return nil, err
	}

	if sourceAccountId == destinationAccountId {
		return nil, fmt.Errorf("cannot create a transaction from and to the same account")
	}

	if sourceAccountId != -1 {
		sourceAccountBalance, err := GetAccountBalance(db, sourceAccountId)
		if err != nil {
			return nil, err
		}

		if sourceAccountBalance < amount {
			return nil, fmt.Errorf("account %d does not have enough avaialble balance (%f) for the transaction (%f)", sourceAccountId, sourceAccountBalance, amount)
		}
	}

	records, err = db.Query("INSERT INTO account_transaction (source_account_id, destination_account_id, amount) VALUES (?, ?, ?) RETURNING *", sourceAccountId, destinationAccountId, amount)
	if err != nil {
		return nil, err
	}

	transactions, err := ScanTransactions(records)
	if err != nil {
		return nil, err
	}

	if len(transactions) != 1 {
		return nil, fmt.Errorf("could not insert transaction!")
	}

	return &transactions[0], nil
}
