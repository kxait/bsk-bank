package handlers

import (
	"bsk-bank/lib"
	"bsk-bank/views"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardConfigHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		records, err := db.Query("SELECT * FROM config")
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		config, err := lib.ScanConfig(records)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		ctx.HTML(http.StatusOK, "", views.Config(config))
	}
}

func PostDashboardConfigSet(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		value := ctx.PostForm("value")
		key := ctx.PostForm("key")

		if value == "" || key == "" {
			if err != nil {
				views.ErrorPage(ctx, "body fields value and key are required")
				return
			}
		}

		records, err := db.Query("UPDATE config SET value = ?, modified_at = CURRENT_TIMESTAMP WHERE key = ? RETURNING *", value, key)

		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		config, err := lib.ScanConfig(records)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		if len(config) != 1 {
			views.ErrorPage(ctx, "Nie udalo sie zmienic konfiguracji :(")
			return
		}

		ctx.Redirect(http.StatusFound, "/dashboard/config")
	}
}
