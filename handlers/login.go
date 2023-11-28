package handlers

import (
	"bsk-bank/lib"
	"bsk-bank/views"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")

		records, err := db.Query("SELECT * FROM user WHERE username = ? AND deleted = 0", username)
		if err != nil {
			lib.ErrorPage(ctx, err.Error())
			return
		}

		users, err := lib.MapUsers(records)
		if err != nil {
			lib.ErrorPage(ctx, err.Error())
			return
		}

		if len(users) != 1 {
			lib.ErrorPage(ctx, "Authentication error! "+fmt.Sprintf("%d", len(users)))
			return
		}

		user := users[0]

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			lib.ErrorPage(ctx, "Authentication error!")
			return
		}

		ctx.HTML(http.StatusOK, "", views.AfterLogin(username, password))
	}
}
