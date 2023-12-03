package lib

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUser(ctx *gin.Context, db *sql.DB) (*User, error) {
	token, err := ctx.Cookie("token")
	if err != nil {
		return nil, err
	}

	records, err := db.Query("SELECT * FROM user WHERE id IN (SELECT user_id FROM session WHERE token = ? AND valid = 1)", token)
	if err != nil {
		return nil, err
	}

	users, err := ScanUsers(records)
	if err != nil {
		return nil, err
	}

	if len(users) != 1 {
		return nil, fmt.Errorf("Could not get user by token! amount of users returned: %d", len(users))
	}

	return &users[0], nil
}
