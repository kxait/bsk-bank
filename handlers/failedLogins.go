package handlers

import (
	"bsk-bank/lib"
	"bsk-bank/views"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardFailedLogsHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		records, err := db.Query("SELECT * FROM failed_login ORDER BY id DESC")
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		logins, err := lib.ScanFailedLogins(records)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		ctx.HTML(http.StatusOK, "", views.DashboardLogins(logins))
	}
}
