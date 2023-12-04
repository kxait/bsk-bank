package lib

import (
	"database/sql"
	"time"
)

type User struct {
	Id       int64
	Username string
	Password string
	Deleted  bool
}

func ScanUsers(rows *sql.Rows) ([]User, error) {
	got := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Deleted)
		if err != nil {
			return []User{}, err
		}
		got = append(got, user)
	}
	return got, nil
}

type Session struct {
	Id        int64
	Token     string
	UserId    int64
	CreatedAt time.Time
	ExpiresAt time.Time
	Valid     bool
}

func ScanSessions(rows *sql.Rows) ([]Session, error) {
	got := make([]Session, 0)
	for rows.Next() {
		var session Session
		var createdAt string
		var expiresAt string
		var valid int64
		err := rows.Scan(&session.Id, &session.Token, &session.UserId, &createdAt, &expiresAt, &valid)
		if err != nil {
			return []Session{}, err
		}

		session.CreatedAt, err = time.Parse(time.DateTime, createdAt)
		if err != nil {
			return []Session{}, err
		}

		session.ExpiresAt, err = time.Parse(time.DateTime, expiresAt)
		if err != nil {
			return []Session{}, err
		}

		if valid == 0 {
			session.Valid = false
		} else {
			session.Valid = true
		}

		got = append(got, session)
	}
	return got, nil
}

type Account struct {
	Id         int64
	HolderName string
	CreatedAt  time.Time
	Deleted    bool
}

func ScanAccounts(rows *sql.Rows) ([]Account, error) {
	got := make([]Account, 0)
	for rows.Next() {
		var account Account
		var createdAt string
		var deleted int64

		err := rows.Scan(&account.Id, &account.HolderName, &createdAt, &deleted)
		if err != nil {
			return []Account{}, err
		}

		account.CreatedAt, err = time.Parse(time.DateTime, createdAt)
		if err != nil {
			return []Account{}, err
		}

		if deleted == 0 {
			account.Deleted = false
		} else {
			account.Deleted = true
		}

		got = append(got, account)

	}
	return got, nil
}

type Transaction struct {
	Id                   int64
	SourceAccountId      int64
	DestinationAccountId int64
	Amount               float64
}

func ScanTransactions(rows *sql.Rows) ([]Transaction, error) {
	got := make([]Transaction, 0)
	for rows.Next() {
		var transaction Transaction

		err := rows.Scan(&transaction.Id, &transaction.SourceAccountId, &transaction.DestinationAccountId, &transaction.Amount)
		if err != nil {
			return []Transaction{}, err
		}

		got = append(got, transaction)
	}

	return got, nil
}

type FailedLogin struct {
	Id        int64
	Username  string
	CreatedAt time.Time
	IpAddress string
}

func ScanFailedLogins(rows *sql.Rows) ([]FailedLogin, error) {
	got := make([]FailedLogin, 0)
	for rows.Next() {
		var failedLogin FailedLogin
		var createdAt string

		err := rows.Scan(&failedLogin.Id, &failedLogin.Username, &createdAt, &failedLogin.IpAddress)
		if err != nil {
			return []FailedLogin{}, err
		}

		failedLogin.CreatedAt, err = time.Parse(time.DateTime, createdAt)
		if err != nil {
			return []FailedLogin{}, err
		}

		got = append(got, failedLogin)
	}

	return got, nil
}

type Config struct {
	Id         int64
	Key        string
	Value      string
	ModifiedAt *time.Time
}

func ScanConfig(rows *sql.Rows) ([]Config, error) {
	got := make([]Config, 0)
	for rows.Next() {
		var config Config

		var modifiedAt sql.NullString
		err := rows.Scan(&config.Id, &config.Key, &config.Value, &modifiedAt)
		if err != nil {
			return []Config{}, err
		}

		if modifiedAt.Valid {
			stuff, err := time.Parse(time.DateTime, modifiedAt.String)
			if err != nil {
				return []Config{}, err
			}
			config.ModifiedAt = &stuff
		}

		got = append(got, config)
	}

	return got, nil
}
