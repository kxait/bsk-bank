package handlers

import (
	"bsk-bank/lib"
	"bsk-bank/views"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		alreadyUser, err := lib.GetUser(ctx, db)
		fmt.Printf("%+v, %s\n", alreadyUser, err)
		if err == nil && alreadyUser != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/dashboard")
			return
		}

		ctx.HTML(http.StatusOK, "", views.Hello())
	}
}
