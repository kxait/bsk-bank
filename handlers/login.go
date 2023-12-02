package handlers

import (
	"bsk-bank/lib"
	"bsk-bank/views"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetLoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		alreadyUser, err := lib.GetUser(ctx, db)
		if err == nil || alreadyUser != nil {
			lib.ErrorPage(ctx, "Log out first")
			return
		}

		ctx.HTML(http.StatusOK, "", views.Login())
	}
}

func PostLoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		alreadyUser, err := lib.GetUser(ctx, db)
		if err == nil || alreadyUser != nil {
			lib.ErrorPage(ctx, "Log out first")
			return
		}

		username := ctx.PostForm("username")
		password := ctx.PostForm("password")

		records, err := db.Query("SELECT * FROM user WHERE username = ? AND deleted = 0", username)
		if err != nil {
			lib.ErrorPage(ctx, err.Error())
			return
		}

		users, err := lib.ScanUsers(records)
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

		token := randomToken()
		expires := time.Now().Add(24 * time.Hour)

		_, err = db.Query("UPDATE session SET valid = 0 WHERE user_id = ?", user.Id)
		if err != nil {
			lib.ErrorPage(ctx, err.Error())
			return
		}

		records, err = db.Query("INSERT INTO session (token, user_id, expires_at) VALUES (?, ?, ?) RETURNING *", token, user.Id, expires)
		if err != nil {
			lib.ErrorPage(ctx, err.Error())
			return
		}

		sessions, err := lib.ScanSessions(records)
		if err != nil {
			lib.ErrorPage(ctx, err.Error())
			return
		}

		if len(sessions) != 1 {
			lib.ErrorPage(ctx, "Could not create user session!")
			return
		}
		fmt.Printf("%+v\n", sessions[0])

		ctx.Header("Set-Cookie", "token="+token)
		ctx.Redirect(http.StatusFound, "/dashboard")
	}
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomToken() string {
	charset := "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset)-1)]
	}
	return string(b)
}

func LogoutHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := lib.GetUser(ctx, db)
		if err == nil && user != nil {
			db.Query("UPDATE session SET valid = 0 WHERE user_id = ?", user.Id)
		}
		ctx.Header("Set-Cookie", "token=; Max-Age=-1")
		ctx.Redirect(http.StatusFound, "/")
	}
}
