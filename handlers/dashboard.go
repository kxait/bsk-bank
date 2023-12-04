package handlers

import (
	"bsk-bank/lib"
	"bsk-bank/views"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		ctx.HTML(http.StatusOK, "", views.Dashboard())
	}
}
