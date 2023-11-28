package lib

import "database/sql"

type User struct {
	Id       int64
	Username string
	Password string
	Deleted  bool
}

func MapUsers(rows *sql.Rows) ([]User, error) {
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
