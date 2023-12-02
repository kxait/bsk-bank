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

		session.ExpiresAt, err = time.Parse("2006-01-02 15:04:05.000000-07:00", expiresAt)
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
